import Hls from "hls.js";
import type { IVideoAnalyzer, VideoMetadata } from "../../domain/video/video.service";

// cf. https://github.com/video-dev/hls.js
export class HlsVideoAnalyzer implements IVideoAnalyzer {
  async analyzeFromUrl(url: string): Promise<VideoMetadata> {
    return new Promise((resolve, reject) => {

      // NOTE: メモリやリスナーを独立させるため、毎回新しいHlsインスタンスを作成する
      const hls = new Hls();

      hls.loadSource(url);

      // 成功イベントリスナー
      hls.on(Hls.Events.MANIFEST_PARSED, (_event, data) => {
        const levels = data.levels;
        if (levels.length === 0) {
          hls.destroy();
          return reject(new Error("No video levels found in HLS manifest"));
        }

        const qualities = levels
          .map((l) => `${l.height}p`)
          .filter((q, index, self) => self.indexOf(q) === index); // 重複削除

        const best = levels[0]

        resolve({
          duration: 0, // HLSは全体のdurationがわからないことが多い
          width: best.width,
          height: best.height,
          bitrate: best.bitrate,
          hasAudio: true, // 一旦true固定
          qualities,
        });

        hls.destroy();
      });

      // エラーイベントリスナー
      hls.on(Hls.Events.ERROR, (_event, data) => {
        if (data.fatal) {
          hls.destroy();
          reject(new Error(`HLS analysis failed: ${data.type} - ${data.details}`));
        }
      });
    });
  }

  async analyzeFromFile(file: File): Promise<VideoMetadata> {
    return new Promise((resolve, reject) => {
      // NOTE: ローカルファイルは、一度ブラウザに読ませる
      const video = document.createElement("video");
      video.preload = "metadata";

      const objectUrl = URL.createObjectURL(file);
      video.src = objectUrl;

      video.onloadedmetadata = () => {
        URL.revokeObjectURL(objectUrl); // メモリ解放

        const duration = video.duration;
        const width = video.videoWidth;
        const height = video.videoHeight;

        // Clean up video element to allow proper garbage collection
        video.onloadedmetadata = null;
        video.onerror = null;
        video.src = "";

        resolve({
          duration,
          width,
          height,
          bitrate: 0, // local からはビットレートはわからないため0固定
          hasAudio: true, // 一旦true固定
          qualities: [`${height}p`],
        });
      };

      video.onerror = (e) => {
        URL.revokeObjectURL(objectUrl); // メモリ解放

        // Clean up video element to allow proper garbage collection
        video.onloadedmetadata = null;
        video.onerror = null;
        video.src = "";

        reject(new Error(`Video analysis failed: ${e}`));
      };
    });
  }
}