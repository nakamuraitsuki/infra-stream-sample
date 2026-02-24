export const hashStringToIndex = (str: string, max: number): number => {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = (hash * 31 + str.charCodeAt(i)) >>> 0; // 符号なし右シフトで32ビット整数に変換
  }
  return hash % max;
}