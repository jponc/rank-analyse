package types

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jponc/rank-analyse/pkg/textrazor"
)

type AnalyzeEntityArray []AnalyzeEntity

type AnalyzeEntity struct {
	ID              uuid.UUID `json:"id" db:"id"`
	ResultID        uuid.UUID `json:"result_id" db:"result_id"`
	Entity          string    `json:"entity" db:"entity"`
	ConfidenceScore float32   `json:"confidence_score" db:"confidence_score"`
	RelevanceScore  float32   `json:"relevance_score" db:"relevance_score"`
	MatchedText     string    `json:"matched_text" db:"matched_text"`
}

func (t *AnalyzeEntityArray) Unmarshal(src interface{}) error {
	switch src.(type) {
	case *textrazor.EntityArray:
		return t.unmarshalTextRazorEntityArray(*src.(*textrazor.EntityArray))
	case textrazor.EntityArray:
		return t.unmarshalTextRazorEntityArray(src.(textrazor.EntityArray))
	}

	return fmt.Errorf("Failed to unmarshal types.AnalyzeEntityArray: '%v'", src)
}

func (t *AnalyzeEntityArray) unmarshalTextRazorEntityArray(textrazorEntities textrazor.EntityArray) error {
	entities := AnalyzeEntityArray{}

	for _, trEntity := range textrazorEntities {
		entities = append(entities, AnalyzeEntity{
			Entity:          trEntity.EntityID,
			ConfidenceScore: trEntity.ConfidenceScore,
			RelevanceScore:  trEntity.RelevanceScore,
			MatchedText:     trEntity.MatchedText,
		})
	}

	*t = entities
	return nil
}
