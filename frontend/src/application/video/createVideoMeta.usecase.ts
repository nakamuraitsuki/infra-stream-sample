import type { AuthSession } from "../../domain/auth/auth.model";
import type { VideoId, VideoTag } from "../../domain/video/video.model";
import type { IVideoRepository } from "../../domain/video/video.repository";

export type CreateVideoMetaResult =
  | { type: "success"; videoId: VideoId }
  | { type: "create_failed"; error: string }
  | { type: "unauthenticated"; error: string };

export type CreateVideoMetaDeps = {
  session: AuthSession | null;
  videoRepo: IVideoRepository;
}

export type CreateVideoMetaParams = {
  title: string;
  description: string;
  tags: VideoTag[];
};

export interface ICreateVideoMetaUseCase {
  execute(params: CreateVideoMetaParams): Promise<CreateVideoMetaResult>;
}

export const createVideoMeta =
  ({ session, videoRepo }: CreateVideoMetaDeps): ICreateVideoMetaUseCase => ({
    execute: async ({ title, description, tags }) => {
      if (!session) {
        return { type: "unauthenticated", error: "User is not authenticated" };
      }

      const res = await videoRepo.create(
        title,
        description,
        tags
      );

      if (!res.success) {
        return { type: "create_failed", error: res.error };
      }

      const videoId = res.data.id;

      return { type: "success", videoId };
    }
  });