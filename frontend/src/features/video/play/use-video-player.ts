import { useCallback, useMemo, useState } from "react";
import { useServices } from "@/context/ServiceContext"
import { getPlaybackDetail, type GetPlaybackDetailResult } from "@/application/video/getPlaybackDetail.usecase";
import type { VideoId } from "@/domain/video/video.model";

export const useVideoPlayer = () => {
  const { videoRepo, videoAnalyzer } = useServices();

  const [loading, setLoading] = useState<boolean>(false);

  const execute = useMemo(
    () => getPlaybackDetail({ videoRepo, videoAnalyzer }),
    [videoRepo, videoAnalyzer]
  );

  // NOTE: 動画情報は変更が少ない前提でキャッチする
  const load = useCallback(async (videoId: VideoId): Promise<GetPlaybackDetailResult> => {
    setLoading(true);
    const result = await execute.execute(videoId);
    setLoading(false);
    return result;
  }, [execute]);

  return { load, loading };
}