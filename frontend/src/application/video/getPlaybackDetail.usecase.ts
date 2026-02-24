import type { Video } from "../../domain/video/video.model";
import type { IVideoRepository } from "../../domain/video/video.repository";
import type { IVideoAnalyzer, VideoMetadata } from "../../domain/video/video.service";

export type PlaybackDetail = {
  url: string;
  mimeType: string;
  meta: VideoMetadata;
  video: Video;
}

export type GetPlaybackDetailResult =
  | { type: "success"; detail: PlaybackDetail }
  | { type: "failed"; error: string };

export type GetPlaybackDetailDeps = {
  videoRepo: IVideoRepository;
  videoAnalyzer: IVideoAnalyzer;
};

export interface IGetPlaybackDetailUseCase {
  execute(videoId: string): Promise<GetPlaybackDetailResult>;
}

export const getPlaybackDetail =
  ({ videoRepo, videoAnalyzer }: GetPlaybackDetailDeps): IGetPlaybackDetailUseCase => ({
    execute: async (videoId: string) => {
      try {
        const videoResult = await videoRepo.findByID(videoId);
        if (!videoResult.success) {
          throw new Error(`Video not found: ${videoId}`);
        }
        const playbackInfoResult = await videoRepo.getPlaybackInfo(videoId);
        if (!playbackInfoResult.success) {
          throw new Error(`Failed to get playback info for video: ${videoId}`);
        }
        const metadataResult = await videoAnalyzer.analyzeFromUrl(playbackInfoResult.data.playbackUrl);

        return {
          type: "success",
          detail: {
            url: playbackInfoResult.data.playbackUrl,
            mimeType: playbackInfoResult.data.mimeType,
            meta: metadataResult,
            video: videoResult.data,
          },
        };
      } catch (error) {
        console.error("Error in getPlaybackDetail use case:", error);
        return {
          type: "failed",
          error: error instanceof Error ? error.message : "Unknown error",
        }
      }
    }
  });