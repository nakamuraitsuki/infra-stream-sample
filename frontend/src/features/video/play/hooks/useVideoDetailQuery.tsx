import { getPlaybackDetail } from "@/application/video/getPlaybackDetail.usecase";
import { useServices } from "@/context/ServiceContext"
import type { VideoId } from "@/domain/video/video.model";
import { useQuery } from "@tanstack/react-query";

export const useVideoDetailQuery = (videoId: VideoId) => {
  const { videoRepo, videoAnalyzer } = useServices();

  return useQuery({
    queryKey: ["video-detail", videoId],
    queryFn: async () => {
      const execute = getPlaybackDetail({ videoRepo, videoAnalyzer });
      const result = await execute.execute(videoId);

      if (result.type !== "success") {
        throw new Error(result.error);
      }
      return result.detail;
    },
    // キャッシュなどの設定はここ
  });
};