import type { VideoStatus } from "@/domain/video/video.model"
import styles from "./VideoStatusChip.module.css";

type VideoStatusChipProps = {
  status: VideoStatus;
  isPublic: boolean;
};

export const VideoStatusChip = ({ status, isPublic }: VideoStatusChipProps) => {
  const getLabel = () => {
    if (status === 'ready') return isPublic ? 'Public' : 'Private';
    if (status === 'failed') return 'Failed';
    return 'Processing';
  };

  const getStyle = () => {
    if (status === 'ready') return isPublic ? styles.readyPublic : styles.readyPrivate;
    if (status === 'failed') return styles.failed;
    return styles.processing;
  };

  const chipStyle = getStyle();
  const label = getLabel();

  return (
    <span className={`${styles.chip} ${chipStyle}`}>
      {label}
    </span>
  )
}