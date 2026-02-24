import { failure, success, type Result } from "../../domain/core/result";
import type { UserId } from "../../domain/user/user.model";
import type { VideoId, Video, VideoTag } from "../../domain/video/video.model";
import type { GetPlaybackInfoResponse, IVideoRepository, VideoError } from "../../domain/video/video.repository";

export class VideoRepositoryMock implements IVideoRepository {
  private DEMO_PLAYBACK_INFO: GetPlaybackInfoResponse = {
    playbackUrl: "https://bitdash-a.akamaihd.net/content/sintel/hls/playlist.m3u8",
    mimeType: "application/vnd.apple.mpegurl",
  };
  private mockVideos: Video[] = [
    {
      id: "1",
      title: "My Public Video",
      description: "This is a mock video for testing purposes.",
      tags: ['competition_programming'],
      createdAt: new Date(),
      ownerId: "1" as UserId,
      failureReason: "None",
      visibility: "public",
      status: "ready",
    },
    {
      id: "2",
      title: "My Private Video",
      description: "This is another mock video for testing purposes.",
      tags: ['web_development'],
      createdAt: new Date(),
      ownerId: "1" as UserId,
      failureReason: "None",
      visibility: "private",
      status: "ready",
    },
    {
      id: "3",
      title: "My processing Video",
      description: "This is another mock video for testing purposes.",
      tags: ['infrastructure'],
      createdAt: new Date(),
      ownerId: "1" as UserId,
      failureReason: "None",
      visibility: "public",
      status: "processing",
    },
    {
      id: "4",
      title: "My failed Video",
      description: "This is another mock video for testing purposes.",
      tags: ['game_development'],
      createdAt: new Date(),
      ownerId: "1" as UserId,
      failureReason: "Encoding error",
      visibility: "public",
      status: "failed",
    },
    {
      id: "5",
      title: "Another Public Video",
      description: "This is yet another mock video for testing purposes.",
      tags: ['machine_learning'],
      createdAt: new Date(),
      ownerId: "2" as UserId,
      failureReason: "None",
      visibility: "public",
      status: "ready",
    }
  ];

  async findByID(id: VideoId): Promise<Result<Video, VideoError>> {
    const video = this.mockVideos.find(video => video.id === id);
    if (video) {
      return success(video);
    } else {
      return failure("NOT_FOUND");
    }
  }

  async findPublicVideos(limit: number): Promise<Video[]> {
    return this.mockVideos.filter(video => video.visibility === "public" && video.status === "ready").slice(0, limit);
  }

  async findByTag(tag: VideoTag): Promise<Video[]> {
    return this.mockVideos.filter(video => video.tags.includes(tag));
  }

  async findMyVideos(limit: number): Promise<Video[]> {
    await new Promise((resolve) => setTimeout(resolve, 1500)); // 擬似的な遅延
    return this.mockVideos.filter(video => video.ownerId === "1").slice(0, limit);
  }

  async getPlaybackInfo(id: VideoId): Promise<Result<GetPlaybackInfoResponse, VideoError>> {
    // NOTE: 簡易的に全ての動画に同じ再生URLを返す
    const video = this.mockVideos.find(video => video.id === id);
    if (!video) {
      return failure("NOT_FOUND");
    }
    if (video.status !== "ready") {
      return failure("VALIDATION_ERROR");
    }
    return success(this.DEMO_PLAYBACK_INFO);
  }

  async create(
    _title: string,
    _description: string,
    _tags: VideoTag[],
  ): Promise<Result<Video, VideoError>> {
    return failure("UNKNOWN_ERROR"); // 動画作成はモックではサポートしない
  }

  async uploadSource(_id: VideoId, _file: File, _onProgress: (progress: number) => void): Promise<Result<void, VideoError>> {
    return failure("UNKNOWN_ERROR"); // 動画アップロードはモックではサポートしない
  }
}