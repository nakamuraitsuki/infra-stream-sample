import { z } from "zod";

const VideoStatusSchema = z.enum([
  'initial',
  'uploaded',
  'queued',
  'processing',
  'ready',
  'failed',
]);

const TagSchema = z.enum([
  'competition_programming',
  'web_development',
  'machine_learning',
  'game_development',
  'infrastructure',
  'other',
]);

const VideoVisibilitySchema = z.enum([
  'public',
  'private',
]);

const ListPublicResponseItemSchema = z.object({
  id: z.string(),
  ownerId: z.string(),
  title: z.string(),
  description: z.string(),
  tags: z.array(TagSchema),
  createdAt: z.string(), // ISO8601形式の文字列
});
const ListPublicResponseSchema = z.object({
  items: z.array(ListPublicResponseItemSchema),
});
export type ListPublicResponseDTO = z.infer<typeof ListPublicResponseSchema>;
export const parseListPublicResponse = (data: unknown): ListPublicResponseDTO => {
  return ListPublicResponseSchema.parse(data);
};

const FindByTagItemSchema = z.object({
  id: z.string(),
  ownerId: z.string(),
  title: z.string(),
  description: z.string(),
  tags: z.array(TagSchema),
  createdAt: z.string(), // ISO8601形式の文字列
});
const SearchByTagResponseSchema = z.object({
  items: z.array(FindByTagItemSchema),
});
export type FindByTagResponseDTO = z.infer<typeof SearchByTagResponseSchema>;
export const parseFindByTagResponse = (data: unknown): FindByTagResponseDTO => {
  return SearchByTagResponseSchema.parse(data);
}

const FindMyVideosItemSchema = z.object({
  id: z.string(),
  title: z.string(),
  description: z.string(),
  tags: z.array(TagSchema),
  status: VideoStatusSchema,
  visibility: VideoVisibilitySchema,
  createdAt: z.string(), // ISO8601形式の文字列
});
const FindMyVideosResponseSchema = z.object({
  items: z.array(FindMyVideosItemSchema),
});
export type FindMyVideosResponseDTO = z.infer<typeof FindMyVideosResponseSchema>;
export const parseFindMyVideosResponse = (data: unknown): FindMyVideosResponseDTO => {
  return FindMyVideosResponseSchema.parse(data);
}