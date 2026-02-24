import type { UserId } from "../user/user.model";

export type VideoId = string;

export type VideoStatus =
  | 'initial'
  | 'uploaded'
  | 'queued'
  | 'processing'
  | 'ready'
  | 'failed';

export type VideoTag =
  | 'competition_programming'
  | 'web_development'
  | 'machine_learning'
  | 'game_development'
  | 'infrastructure'
  | 'other';

export type VideoVisibility = 'public' | 'private';

export interface Video {
  readonly id: VideoId;
  readonly title: string;
  readonly description: string;
  readonly tags: VideoTag[];
  readonly createdAt: Date;
  // 状況によってバックエンドから提供されない詳細情報
  readonly ownerId?: UserId;
  readonly failureReason?: string;
  readonly visibility?: VideoVisibility;
  readonly status?: VideoStatus;
}