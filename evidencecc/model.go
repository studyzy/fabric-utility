package main
type DataEvidence struct{
	Owner string  `json:"owner"`
	Category string  `json:"category"`
	DataKey string  `json:"dataKey"`
	DataValue string  `json:"dataValue"`
	Reference string  `json:"reference"`
	Timestamp int64  `json:"timestamp"`
	DocType string `json:"docType"`
}
type ObjectEvidence struct{
	Owner string  `json:"owner"`
	Category string  `json:"category"`
	ObjectKey string  `json:"objectKey"`
	Object map[string]interface{}  `json:"object"`
	Reference string  `json:"reference"`
	Timestamp int64  `json:"timestamp"`
	DocType string `json:"docType"`

}
