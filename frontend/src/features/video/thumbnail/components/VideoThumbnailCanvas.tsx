import type { VideoId, VideoStatus } from "@/domain/video/video.model"
import styles from "./VideoThumbnailCanvas.module.css";
import { useEffect, useRef, useState } from "react";

type Props = {
  videoId: VideoId;
  title: string;
  status?: VideoStatus;
}

export const VideoThumbnailCanvas = ({ title, status = "ready" }: Props) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [src, setSrc] = useState<string>("");

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    const img = new Image();
    img.src = `/thumbnail/cover1.jpg`;
    img.onload = () => {
      canvas.width = img.width;
      canvas.height = img.height;

      ctx.drawImage(img, 0, 0);

      if (status !== "ready") {
        ctx.fillStyle = "rgba(0, 0, 0, 0.5)";
        ctx.fillRect(0, 0, canvas.width, canvas.height);
      }

      const maxHeight = canvas.height / 8; // 最大8行分の高さ
      const fontSize = Math.min(canvas.width / 6, maxHeight / 1.2);

      ctx.font = `bold ${fontSize}px system-ui`;
      ctx.fillStyle = "black";
      ctx.textAlign = "center";
      ctx.textBaseline = "middle";

      const words = title.split(" ");
      const lineHeight = fontSize * 1.2;
      const maxLines = 4;
      let lines: string[] = [];
      let currentLine = "";

      words.forEach((word) => {
        const testLine = currentLine ? `${currentLine} ${word}` : word;
        const metrics = ctx.measureText(testLine);
        if (metrics.width > canvas.width * 0.8 && currentLine) {
          lines.push(currentLine);
          currentLine = word;
        } else {
          currentLine = testLine;
        }
      });
      if (currentLine) lines.push(currentLine);
      if (lines.length > maxLines) lines = lines.slice(0, maxLines);

      lines.forEach((line, i) => {
        ctx.fillText(line, canvas.width / 2, canvas.height / 2 + (i - (lines.length - 1) / 2) * lineHeight);
      });

      // DataURLとしてセット
      setSrc(canvas.toDataURL());
    };
  }, [title, status]);

  return (
    <div className={styles.thumbnail}>
      {src ? (
        <img src={src} alt={title} className={styles.image} />
      ) : (
        <canvas ref={canvasRef} style={{ display: "none" }} />
      )}
    </div>
  )
}