package eventschema

const (
	ProcessKeyword string = "ProcessKeyword"
	ResultCreated  string = "ResultCreated"
)

type ProcessKeywordMessage struct {
	Keyword      string `json:"keyword"`
	SearchEngine string `json:"search_engine"`
	Device       string `json:"device"`
	Count        int    `json:"count"`
}

type ResultCreatedMessage struct {
	ResultID string `json:"result_id"`
}
