import { getPublicVideos } from "@/application/video/getPublicVideos.usecase";
import { useAuth } from "@/context/AuthContext"
import { useServices } from "@/context/ServiceContext";
import type { Video } from "@/domain/video/video.model";
import { useQuery } from "@tanstack/react-query";

export const usePublicVideosQuery = (limit: number) => {
  const { session } = useAuth();
  const { videoRepo } = useServices();

  return useQuery<{ videos: Video[]}, Error>({
    queryKey: ["publicVideos", limit],
    queryFn: async () => {
      const usecase = getPublicVideos({ videoRepo, session });
      const result = await usecase.execute(limit);

      if (result.type !== "success") {
        throw new Error("Failed to fetch public videos");
      }

      return { videos: result.videos };
    }
    // キャッシュ設定などを追記
  });
};