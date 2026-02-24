import type { AuthSession } from "../../domain/auth/auth.model";
import type { Video } from "../../domain/video/video.model";
import type { IVideoRepository } from "../../domain/video/video.repository";

export type GetMyVideosResult =
  | { type: "success"; videos: Video[] }
  | { type: "unauthenticated" }
  | { type: "failed"; error: string };

export type GetMyVideosDeps = {
  videoRepo: IVideoRepository;
  session: AuthSession | null;
};

export interface IGetMyVideosUseCase {
  execute(limit: number): Promise<GetMyVideosResult>;
}

export const getMyVideos =
  ({ videoRepo, session }: GetMyVideosDeps): IGetMyVideosUseCase => ({
    execute: async (limit: number) => {
      if (!session) {
        return { type: "unauthenticated" };
      }
      try {
        const videos = await videoRepo.findMyVideos(limit);
        return { type: "success", videos };
      } catch (error) {
        return { type: "failed", error: (error as Error).message };
      }
    }
  });