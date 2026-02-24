import { useForm } from "react-hook-form";
import { useCreateVideoMutation } from "../../hooks/useCreateVideoMutation";
import styles from "./VideoUploadForm.module.css";
import { IconButton } from "@/ui/IconButton/IconButton";
import { BiArrowBack } from "react-icons/bi";
import { AiOutlineUpload } from "react-icons/ai";
import { useEffect } from "react";

type FormValues = {
  title: string;
  description: string;
  videoFile: FileList;
}

type Props = {
  onSuccess: (videoId: string) => void;
  onBack: () => void;
}

export const VideoUploadForm = ({ onSuccess, onBack }: Props) => {
  const { mutate, isPending, progress, error: apiError } = useCreateVideoMutation();

  const {
    register,
    handleSubmit,
    formState: { errors, isValid }
  } = useForm<FormValues>({
    mode: "onChange",
  });

  // よくあるページ移動に対する警告
  useEffect(() => {
    if (!isPending) return;

    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      e.preventDefault();
      return ""; // NOTE: 古っるいブラウザのためのやつ
    };

    window.addEventListener("beforeunload", handleBeforeUnload);
    return () => window.removeEventListener("beforeunload", handleBeforeUnload);
  }, [isPending]);

  const onSubmit = (data: FormValues) => {
    const file = data.videoFile[0];
    if (!file) return;

    mutate(
      {
        title: data.title,
        description: data.description,
        tags: [],
        file,
      },
      { onSuccess: (res) => onSuccess(res.videoId) }
    );
  };

  return (
    <div className={styles.container}>
      <div className={styles.nav}>
        <IconButton
          icon={<BiArrowBack />}
          onClick={onBack}
          label="Back"
          disabled={isPending}
        />
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className={styles.form}>
        <div className={styles.field}>
          <label htmlFor="title">Title</label>
          <input
            id="title"
            {...register("title", {
              required: "Title is required",
              maxLength: { value: 100, message: "Max length is 100 characters" },
            })}
            disabled={isPending}
          />
          {errors.title && <p className={styles.error}>{errors.title.message}</p>}
        </div>

        <div className={styles.field}>
          <label htmlFor="description">Description</label>
          <textarea
            id="description"
            {...register("description")}
            className={styles.textarea}
            disabled={isPending}
          />
        </div>

        <div className={styles.field}>
          <label htmlFor="videoFile">Video File</label>
          <input
            id="videoFile"
            type="file"
            accept="video/*"
            {...register("videoFile", {
              required: "Video file is required"
            })}
            disabled={isPending}
            className={styles.fileInput}
          />
          {errors.videoFile && <p className={styles.error}>{errors.videoFile.message}</p>}
        </div>

        {isPending && (
          <div className={styles.progressSection}>
            <div className={styles.progressBarContainer}>
              <div
                className={styles.progressBar}
                style={{ width: `${progress}%` }}
              />
            </div>
            <p className={styles.progressLabel}>
              {progress < 100 ? `Uploading... ${progress.toFixed(2)}%` : "Finalizing..."}
            </p>
          </div>
        )}

        {apiError && <p className={styles.error}>Error: {apiError.message}</p>}

        <IconButton
          icon={<AiOutlineUpload size={24} />}
          type="submit"
          disabled={!isValid || isPending}
          label="Upload Video"
        />
      </form >
    </div >
  )
}
