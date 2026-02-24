import type { Result } from "../core/result";
import type { VideoTag, Video, VideoId } from "./video.model";

export interface GetPlaybackInfoResponse {
  playbackUrl: string; // Master playlist URL (e.g., HLS)
  mimeType: string;
}

export interface IVideoRepository {
  findByID(id: VideoId): Promise<Result<Video, VideoError>>;

  findPublicVideos(limit: number): Promise<Video[]>;
  findByTag(tag: VideoTag): Promise<Video[]>;
  findMyVideos(limit: number): Promise<Video[]>;

  getPlaybackInfo(id: VideoId): Promise<Result<GetPlaybackInfoResponse, VideoError>>;
  create(
    title: string,
    description: string,
    tags: VideoTag[],
  ): Promise<Result<Video, VideoError>>;
  uploadSource(id: VideoId, file: File, onProgress?: (progress: number) => void): Promise<Result<void, VideoError>>;
}

export type VideoError = 
  | "NOT_FOUND"
  | "UNAUTHORIZED"
  | "VALIDATION_ERROR"
  | "SERVER_ERROR"
  | "UNKNOWN_ERROR";