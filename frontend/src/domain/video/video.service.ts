export interface VideoMetadata {
  duration: number; // in seconds
  width: number;
  height: number;
  bitrate: number;
  hasAudio: boolean;
  qualities: string[]; // e.g., ['1080p', '720p']
}

export interface IVideoAnalyzer {
  analyzeFromUrl(url: string): Promise<VideoMetadata>;

  // バリデーション用　アップ前の動画解析
  analyzeFromFile(file: File): Promise<VideoMetadata>;
}