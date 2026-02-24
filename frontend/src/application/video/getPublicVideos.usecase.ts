import type { AuthSession } from "../../domain/auth/auth.model";
import type { Video } from "../../domain/video/video.model";
import type { IVideoRepository } from "../../domain/video/video.repository";

export type GetPublicVideosResult =
  | { type: "success"; videos: Video[] }
  | { type: "failed"; error: string };

export type GetPublicVideosDeps = {
  videoRepo: IVideoRepository;
  session: AuthSession | null;
};

export interface IGetPublicVideosUseCase {
  execute(limit: number): Promise<GetPublicVideosResult>;
}

export const getPublicVideos =
  ({ videoRepo }: GetPublicVideosDeps): IGetPublicVideosUseCase => ({
    execute: async (limit: number) => {
      try {
        const videos = await videoRepo.findPublicVideos(limit);
        return { type: "success", videos };
      } catch (error) {
        return { type: "failed", error: (error as Error).message };
      }
    }
  });