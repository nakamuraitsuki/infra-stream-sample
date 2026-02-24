package process

import "github.com/google/uuid"

// DTO for transcode job payload
type TranscodePayload struct {
	VideoID uuid.UUID `json:"video_id"`
}
