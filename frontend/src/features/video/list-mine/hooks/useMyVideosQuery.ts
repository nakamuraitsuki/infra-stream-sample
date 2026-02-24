import { useAuth } from "@/context/AuthContext";
import { useServices } from "@/context/ServiceContext";
import { getMyVideos } from "@/application/video/getMyVideos.usecase";
import { useQuery } from "@tanstack/react-query";
import type { Video } from "@/domain/video/video.model";

export const useMyVideosQuery = (limit: number) => {
  const { session } = useAuth();
  const { videoRepo } = useServices();

  const execute = getMyVideos({ videoRepo, session });

  return useQuery<{ videos: Video[] }, Error>({
    queryKey: ["myVideos", limit],
    queryFn: async () => {
      const res = await execute.execute(limit);

      if (res.type !== "success") {
        throw new Error("Failed to fetch my videos");
      }

      return { videos: res.videos };
    },
    enabled: session.status === "authenticated", // ログイン時のみクエリを有効化
    // キャッシュなどが必要な場合は追記
  });
};