import type { PlaybackDetail } from "@/application/video/getPlaybackDetail.usecase";
import Hls from "hls.js";
import { useEffect, useRef } from "react";
import styles from "./VideoPlayerView.module.css";

export const VideoPlayerView = ({ detail }: { detail: PlaybackDetail }) => {
  const videoRef = useRef<HTMLVideoElement>(null);

  // hls初期化とクリーンアップ
  useEffect(() => {
    const video = videoRef.current;
    if (!video || !detail.url) return;

    let hls: Hls;

    if (Hls.isSupported()) {
      hls = new Hls();
      hls.loadSource(detail.url);
      hls.attachMedia(video);
    } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
      video.src = detail.url;
    } else {
      console.error("HLS is not supported in this browser.");
    }

    return () => {
      if (hls) {
        hls.destroy();
      }
    };
  }, [detail.url]);

  return (
    <div className={styles.playerWrapper}>
      <h2 className={styles.title}>{detail.video.title}</h2>

      <div className={styles.videoContainer}>
        <video ref={videoRef} controls className={styles.videoElement} />
      </div>

      <div className={styles.info}>
        <p className={styles.description}>{detail.video.description}</p>
      </div>
    </div>
  )
}