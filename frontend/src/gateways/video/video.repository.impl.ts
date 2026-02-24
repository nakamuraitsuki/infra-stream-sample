import axios from "axios";
import { apiClient } from "../../api/client";
import type { VideoId, Video, VideoTag, VideoStatus, VideoVisibility } from "../../domain/video/video.model";
import type { GetPlaybackInfoResponse, IVideoRepository, VideoError } from "../../domain/video/video.repository";
import { failure, success, type Result } from "../../domain/core/result";
import { parseFindByTagResponse, parseFindMyVideosResponse, parseListPublicResponse } from "./video.dto";
import pLimit from "p-limit";

export class VideoRepositoryImpl implements IVideoRepository {
  // helper method to map axios errors to VideoError
  private handleError(error: unknown): VideoError {
    if (axios.isAxiosError(error)) {
      const status = error.response?.status;
      if (status === 404) return "NOT_FOUND";
      if (status === 401) return "UNAUTHORIZED";
      if (status === 400) return "VALIDATION_ERROR";
      if (status && status >= 500) return "SERVER_ERROR";
    }
    return "UNKNOWN_ERROR";
  }

  async findByID(id: VideoId): Promise<Result<Video, VideoError>> {
    try {
      const { data } = await apiClient.get(`/api/videos/${id}`);
      return success(data);
    } catch (error) {
      return failure(this.handleError(error));
    }
  }

  async findPublicVideos(limit: number): Promise<Video[]> {
    try {
      const { data } = await apiClient.get("/api/videos", {
        params: { limit },
      });

      const dto = parseListPublicResponse(data);
      return dto.items.map((item): Video => ({
        id: item.id,
        ownerId: item.ownerId,
        title: item.title,
        description: item.description,
        tags: item.tags as VideoTag[],
        createdAt: new Date(item.createdAt),
      }))
    } catch (error) {
      console.error("Failed to fetch public videos:", error);
      return [];
    }
  }

  async findByTag(tag: VideoTag): Promise<Video[]> {
    try {
      const { data } = await apiClient.get("/api/videos/search", {
        params: { tag },
      })
      const dto = parseFindByTagResponse(data);
      return dto.items.map((item): Video => ({
        id: item.id,
        ownerId: item.ownerId,
        title: item.title,
        description: item.description,
        tags: item.tags as VideoTag[],
        createdAt: new Date(item.createdAt),
      }));
    } catch (error) {
      console.error(`Failed to search videos by tag "${tag}":`, error);
      return [];
    }
  }

  async findMyVideos(limit: number): Promise<Video[]> {
    try {
      const { data } = await apiClient.get("/api/videos/mine", {
        params: { limit },
      });
      const dto = parseFindMyVideosResponse(data);
      return dto.items.map((item): Video => ({
        id: item.id,
        title: item.title,
        description: item.description,
        tags: item.tags as VideoTag[],
        status: item.status as VideoStatus,
        visibility: item.visibility as VideoVisibility,
        createdAt: new Date(item.createdAt),
      }));
    } catch (error) {
      console.error("Failed to fetch my videos:", error);
      return [];
    }
  }

  async getPlaybackInfo(id: VideoId): Promise<Result<GetPlaybackInfoResponse, VideoError>> {
    try {
      const { data } = await apiClient.get<GetPlaybackInfoResponse>(
        `/api/videos/${id}/playback-info`,
      );
      return success(data);
    } catch (error) {
      return failure(this.handleError(error));
    }
  }

  async create(title: string, description: string, tags: VideoTag[]): Promise<Result<Video, VideoError>> {
    try {
      const { data } = await apiClient.post<Video>("/api/videos", {
        title,
        description,
        tags,
      });
      return success(data);
    } catch (error) {
      return failure(this.handleError(error));
    }
  }

  async uploadSource(id: VideoId, file: File, onProgress: (progress: number) => void = () => {}): Promise<Result<void, VideoError>> {
    const PROMISE_LIMIT = 3;
    try {
      const { data: initData } = await apiClient.post<{
        uploadId: string;
        urls: string[];
        partSize: number;
        key: string;
      }>(`/api/videos/${id}/upload/init`, {
        fileSize: file.size,
      });

      const { uploadId, urls, partSize } = initData;

      // 各パーツの正確な期待サイズを事前計算（進捗反映のため）
      const partExpectedSizes = urls.map((_, index) => {
        const start = index * partSize;
        const end = Math.min(start + partSize, file.size);
        return end - start;
      });

      // 各パーツの現在のアップロード済みバイト数
      const uploadedBytes = new Array(urls.length).fill(0);

      const updateTotalProgress = (forceComplete = false) => {
        if (forceComplete) {
          onProgress(100);
          return;
        }
        const totalUploaded = uploadedBytes.reduce((acc, curr) => acc + curr, 0);
        const progress = Math.min(99, Math.floor((totalUploaded / file.size) * 100));
        onProgress(progress);
      };

      const limit = pLimit(PROMISE_LIMIT);

      const uploadPromises = urls.map(async (url, index) => {
        return limit(async () => {
          const partNumber = index + 1;
          const start = index * partSize;
          const end = Math.min(start + partSize, file.size);
          const chunk = file.slice(start, end);
          const expectedSize = partExpectedSizes[index];

          let lastError: any;
          for (let attempt = 1; attempt <= 3; attempt++) {
            try {
              // リトライ時はパーツの進捗を一旦0にリセットして計算（進捗反映のため）
              uploadedBytes[index] = 0;
              updateTotalProgress();

              const response = await axios.put(url, chunk, {
                headers: { "Content-Type": file.type || "video/mp4" },
                onUploadProgress: (progressEvent) => {
                  // Axiosの進捗が期待サイズを超えないようガードしつつ記録
                  uploadedBytes[index] = Math.min(progressEvent.loaded, expectedSize);
                  updateTotalProgress();
                },
              });

              const etag = response.headers.etag;
              if (!etag) throw new Error("ETag missing");

              // 成功時は確実にそのパーツ分をフルカウント
              uploadedBytes[index] = expectedSize;
              updateTotalProgress();

              return {
                partNumber: partNumber,
                etag: etag.replace(/"/g, ""),
              };
            } catch (err) {
              lastError = err;
              if (attempt < 3) {
                const delayMs = Math.pow(2, attempt - 1) * 1000;
                await new Promise(resolve => setTimeout(resolve, delayMs));
              }
            }
          }
          throw lastError;
        });
      });

      const results = await Promise.all(uploadPromises);

      // サーバー側の結合処理(complete)を呼ぶ前に100%を表示
      updateTotalProgress(true);

      await apiClient.post(`/api/videos/${id}/upload/complete`, {
        uploadId: uploadId,
        parts: results.sort((a, b) => a.partNumber - b.partNumber),
      });

      return success(undefined);
    } catch (error) {
      console.error("Multipart upload failed:", error);
      return failure(this.handleError(error));
    }
  }
}
