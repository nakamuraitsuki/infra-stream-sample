import type { Video, VideoId } from "@/domain/video/video.model"
import styles from "./PublicVideoCard.module.css";
import { VideoThumbnailCanvas } from "@/features/video/thumbnail";

type Props = {
  video: Video;
  onClick?: (videoId: VideoId) => void;
};

export const PublicVideoCard = ({ video, onClick }: Props) => {
  return (
    <div
      className={styles.card}
      onClick={() => onClick?.(video.id)}
    >
      <div className={styles.thumbnail}>
        <VideoThumbnailCanvas
          videoId={video.id}
          title={video.title}
        />
      </div>
      <div className={styles.cardBody}>
        <p className={styles.description}>{video.description}</p>
      </div>
    </div>
  );
};