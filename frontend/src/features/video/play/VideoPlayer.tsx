import { useEffect, useRef, useState } from "react";
import Hls from "hls.js";
import type { PlaybackDetail } from "@/application/video/getPlaybackDetail.usecase";
import { useVideoPlayer } from "./use-video-player";

type Props = {
  videoId: string;
  onErrorBack: () => void;
};

export const VideoPlayer = ({ videoId, onErrorBack }: Props) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const { load, loading } = useVideoPlayer();

  const [detail, setDetail] = useState<PlaybackDetail | null>(null);
  const [error, setError] = useState<string | null>(null);

  // データ取得
  useEffect(() => {
    load(videoId).then((res) => {
      if (res.type === "success") setDetail(res.detail);
      else setError(res.error);
    });
  }, [videoId, load]);

  // HLS初期化
  useEffect(() => {
    const video = videoRef.current;
    if (!video || !detail) return;

    if (!Hls.isSupported()) {
      setError("HLS is not supported in this browser");
      return;
    }

    const hls = new Hls();
    hls.loadSource(detail.url);
    hls.attachMedia(video);

    return () => {
      hls.destroy();
    };
  }, [detail]);

  if (loading) return <div style={{ padding: "20px" }}>Loading video...</div>;

  if (error)
    return (
      <div style={{ padding: "20px", color: "red" }}>
        <p>Error: {error}</p>
        <button onClick={onErrorBack}>Back</button>
      </div>
    );

  return (
    <div style={{ padding: "20px" }}>
      <h2>{detail?.video?.title}</h2>

      <div style={{ maxWidth: "800px", background: "#000" }}>
        <video
          ref={videoRef}
          controls
          style={{ width: "100%", height: "100%" }}
        />
      </div>

      <div style={{ marginTop: "10px" }}>
        <p>{detail?.video?.description}</p>
      </div>
    </div>
  );
};
