import { Spinner } from "@/ui/Spinner/Spinner";
import { useVideoDetailQuery } from "../../hooks/useVideoDetailQuery";
import styles from "./VideoPlayerContainer.module.css";
import { BiArrowBack, BiErrorCircle } from "react-icons/bi";
import { IconButton } from "@/ui/IconButton/IconButton";
import { VideoPlayerView } from "../VideoPlayerView/VideoPlayerView";

type Props = {
  videoId: string;
  onErrorBack: () => void;
}

export const VideoPlayerContainer = ({ videoId, onErrorBack }: Props) => {
  const { data: detail, isLoading, error } = useVideoDetailQuery(videoId);

  if (isLoading) {
    return <Spinner isLoading={isLoading} />;
  }

  if (error || !detail) {
    return (
      <div className={styles.errorContainer}>
        <BiErrorCircle size={48} color="red" />
        <p>Error: {error?.message || "Failed to load video details."}</p>
        <IconButton icon={<BiArrowBack size={24} />} onClick={onErrorBack} label="Back to List" />
      </div>
    );
  }

  return <VideoPlayerView detail={detail} />;
};
