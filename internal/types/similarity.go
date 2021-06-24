package types

type SimilarityKeyword struct {
	Keyword string             `json:"keyword"`
	Results []SimilarityResult `json:"similarity_results"`
}

type SimilarityResult struct {
	Positions []int  `json:"positions"`
	SeenCount int    `json:"seen_count"`
	Title     string `json:"title"`
}
