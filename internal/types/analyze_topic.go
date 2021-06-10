package types

import (
	"fmt"

	"github.com/jponc/rank-analyse/pkg/textrazor"
)

type AnalyzeTopicArray []AnalyzeTopic

type AnalyzeTopic struct {
	Label string  `json:"label" db:"label"`
	Score float32 `json:"score" db:"score"`
}

func (t *AnalyzeTopicArray) Unmarshal(src interface{}) error {
	switch src.(type) {
	case *textrazor.TopicArray:
		return t.unmarshalTextRazorTopicArray(*src.(*textrazor.TopicArray))
	case textrazor.TopicArray:
		return t.unmarshalTextRazorTopicArray(src.(textrazor.TopicArray))
	}

	return fmt.Errorf("Failed to unmarshal types.AnalyzeTopicArray: '%v'", src)
}

func (t *AnalyzeTopicArray) unmarshalTextRazorTopicArray(textrazorTopics textrazor.TopicArray) error {
	topics := AnalyzeTopicArray{}

	for _, trTopic := range textrazorTopics {
		topics = append(topics, AnalyzeTopic{
			Label: trTopic.Label,
			Score: trTopic.Score,
		})
	}

	*t = topics
	return nil
}
