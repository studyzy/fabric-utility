package main
type User struct{
	UserId string `json:"userId"`
	UserName string `json:"userName"`
	OrgName string `json:"orgName"`
	PubKey string `json:"pubKey"`
	DocType string `json:"docType"`

}
type DataEvidence struct{
	Owner string  `json:"owner"`
	Category string  `json:"category"`
	DataKey string  `json:"dataKey"`
	DataValue string  `json:"dataValue"`
	Reference string  `json:"reference"`
	Timestamp int64  `json:"timestamp"`
	UnlockKey string `json:"unlockKey"`
	DocType string `json:"docType"`

}
type DataShare struct{
	Owner string  `json:"owner"`
	Category string  `json:"category"`
	DataKey string  `json:"dataKey"`
	ShareTo string `json:"shareTo"`
	ShareStartTime int64  `json:"shareStartTime"`
	ShareEndTime int64  `json:"shareEndTime"`
	Timestamp int64  `json:"timestamp"`
	UnlockKey string `json:"unlockKey"`
	DocType string `json:"docType"`

}
