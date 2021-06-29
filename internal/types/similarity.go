package types

type SimilarityKeyword struct {
	Keyword string             `json:"keyword"`
	Results []SimilarityResult `json:"similarity_results"`
}

type SimilarityResult struct {
	AvePosition float32 `json:"average_position"`
	SeenCount   int     `json:"seen_count"`
	Title       string  `json:"title"`
	Link        string  `json:"link"`
}
