import type { Video, VideoId } from "@/domain/video/video.model"
import { VideoThumbnailCanvas } from "@/features/video/thumbnail";
import styles from "./MyVideoCard.module.css";
import { VideoStatusChip } from "../../../ui/Chip/VideoStatusChip";

type Props = {
  video: Video;
  onClick?: (videoId: VideoId) => void;
};

export const MyVideoCard = ({ video, onClick }: Props) => {
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
        {video.status ? (
          <VideoStatusChip status={video.status} isPublic={video.visibility === "public"} />
        ) : null}
        <div className={styles.meta}>
          <h3 className={styles.title}>{video.title}</h3>
          <span className={styles.description}>{video.description}</span>
        </div>
      </div>
    </div>
  );
};