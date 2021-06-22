package types

type SimilaryKeyword struct {
	Keyword string
	Results []SimilarityResult
}

type SimilarityResult struct {
	Position  float32
	SeenCount int
	Title     string
}
