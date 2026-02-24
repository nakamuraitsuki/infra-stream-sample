-- 1. Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,     
    bio TEXT NOT NULL DEFAULT '', 
    icon_key TEXT,                  
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- 2. Videos Table
CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_key TEXT NOT NULL,
    stream_key TEXT NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    retry_count INTEGER NOT NULL DEFAULT 0,
    failure_reason TEXT,         
    visibility VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- 3. Tags Master Table
CREATE TABLE IF NOT EXISTS tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- 4. Video-Tags Junction Table
-- videoTagModel の定義 (VideoID: int64, TagID: uuid) に合わせる場合は
-- スキーマ側もその型にする必要がありますが、通常は videos.id(UUID) を参照するため
-- ここでは DTO の意図（リレーション）を汲みつつ videos(id) と tags(id) を結びつけます
CREATE TABLE IF NOT EXISTS video_tags (
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (video_id, tag_id)
);

-- 5. Outbox Table (PayloadをBYTEAに変更)
-- DTO で Payload []byte となっているため、PostgreSQL では BYTEA が最適です
CREATE TABLE IF NOT EXISTS outbox (
    id UUID PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    payload BYTEA NOT NULL,         -- []byte に準拠
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_videos_owner_id ON videos(owner_id);
CREATE INDEX IF NOT EXISTS idx_outbox_occurred_at ON outbox(occurred_at ASC);
CREATE INDEX IF NOT EXISTS idx_video_tags_tag_id ON video_tags(tag_id);

-- テスト用のユーザーを作成
INSERT INTO users (id, name, role, created_at)
VALUES ('00000000-0000-0000-0000-000000000000', 'testuser', 'user', NOW())
ON CONFLICT (id) DO NOTHING;