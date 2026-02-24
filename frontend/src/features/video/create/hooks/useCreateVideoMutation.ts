import { createVideoMeta } from "@/application/video/createVideoMeta.usecase";
import { uploadVideoSource } from "@/application/video/uploadVideoSource.usecase";
import { useAuth } from "@/context/AuthContext";
import { useServices } from "@/context/ServiceContext";
import type { VideoTag } from "@/domain/video/video.model";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";

export const useCreateVideoMutation = () => {
  const { session } = useAuth();
  const { videoRepo } = useServices();

  const [progress, setProgress] = useState(0);

  // cf. https://zenn.dev/taisei_13046/books/133e9995b6aadf/viewer/257b1a
  const mutation = useMutation({
    mutationFn: async ({
      title,
      description,
      tags,
      file,
    }: {
      title: string;
      description: string;
      tags: VideoTag[];
      file: File;
    }) => {
      setProgress(0);

      const createMeta = createVideoMeta({ session, videoRepo });
      const createRes = await createMeta.execute({ title, description, tags });

      if (createRes.type !== "success") {
        throw new Error(`${createRes.type} - ${createRes.error}`);
      }

      const uploadSource = uploadVideoSource({ videoRepo });
      const uploadRes = await uploadSource.execute({
        videoId: createRes.videoId,
        file,
        onProgress: (progress) => setProgress(progress),
      });

      if (uploadRes.type !== "success") {
        throw new Error(`${uploadRes.type} - ${uploadRes.error}`);
      }

      return { videoId: createRes.videoId };
    }
  });

  return { ...mutation, progress };
};