import type { Video, VideoId } from "@/domain/video/video.model";
import { usePublicVideosQuery } from "../../hooks/usePublicVideosQuery";
import { PublicVideoCard } from "../PublicVideoCard/PublicVideoCard";
import styles from "./PublicVideoList.module.css";
import { Spinner } from "@/ui/Spinner/Spinner";

type Props = {
  limit?: number;
  onSelect?: (videoId: VideoId) => void;
}

export const PublicVideoList = ({ limit = 20, onSelect }: Props) => {
  const { data, isLoading, error } = usePublicVideosQuery(limit);

  if (isLoading) return <Spinner isLoading={isLoading} />
  if (error) return <p>{error.message}</p>
  if (!data || data.videos.length === 0) return <p>No public videos available.</p>

  return (
    <div className={styles.grid}>
      {data.videos.map((video: Video) => (
        <PublicVideoCard
          key={video.id}
          video={video}
          onClick={onSelect}
        />
      ))}
    </div>
  );
};