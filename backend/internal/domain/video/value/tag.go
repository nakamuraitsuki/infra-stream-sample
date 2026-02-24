package value

import "errors"

type Tag string

const (
	TagCompetitionProgramming Tag = "competition_programming"
	TagWebDevelopment         Tag = "web_development"
	TagMachineLearning        Tag = "machine_learning"
	TagGameDevelopment        Tag = "game_development"
	TagInfrastructure         Tag = "infrastructure"
)

var (
	ErrInvalidTag = errors.New("invalid tag")
)

func NewTag(s string) (Tag, error) {
	switch s {
	case string(TagCompetitionProgramming),
		string(TagWebDevelopment),
		string(TagMachineLearning),
		string(TagGameDevelopment),
		string(TagInfrastructure):
		return Tag(s), nil
	default:
		return "", ErrInvalidTag
	}
}
