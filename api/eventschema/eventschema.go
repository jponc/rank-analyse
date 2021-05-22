package eventschema

const (
	ProcessKeyword string = "ProcessKeyword"
	ResultCreated  string = "ResultCreated"
	CrawlFinished  string = "CrawlFinished"
)

type ProcessKeywordMessage struct {
	Keyword      string `json:"keyword"`
	SearchEngine string `json:"search_engine"`
	Device       string `json:"device"`
	Count        int    `json:"count"`
	Email        string `json:"email"`
}

type ResultCreatedMessage struct {
	ResultID string `json:"result_id"`
}

type CrawlFinishedMessage struct {
	CrawlID string `json:"crawl_id"`
}
