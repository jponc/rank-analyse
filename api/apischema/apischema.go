package apischema

type HealthcheckResponse struct {
	Status string `json:"status"`
}

type LambdaTestResponse struct {
	Out string `json:"out"`
}
