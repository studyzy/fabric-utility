package main

type CustomerInfo struct {
	Id           int    `json:"id"`         //key
	CustomerId   int    `json:"customerId"` //key
	CustomerName string `json:"customerName"`
	Balance      uint64 `json:"balance"`
	Timestamp    int64  `json:"timestamp"` //最后更新时间
	DocType      string `json:"docType"`   //用来区分不同的Struct的，这里的值是CustomerInfo
}
type UploadData struct {
	Id             int    `json:"id"`         //key
	Category       string `json:"category"`   //上传的数据类型
	CreateTime     int64  `json:"createTime"` //上传时间
	CustomerId     int    `json:"customerId"` //
	CustomerName   string `json:"customerName"`
	Interal        uint64 `json:"interal"` //很奇怪的一个词，这里表示获得的积分
	Hash           string `json:"hash"`
	BatchDataTotal int    `json:"batchDataTotal"`
	Timestamp      int64  `json:"timestamp"` //区块链存证时间
	DocType        string `json:"docType"`   //用来区分不同的Struct的，这里的值是UploadData
}
type ExchangeOrder struct {
	Id             int    `json:"id"`
	OrderNo        string `json:"orderNo"`
	Category       string `json:"category"`
	CreateTime     int64  `json:"createTime"`
	DataProvider   string `json:"dataProvider"`
	DataProviderId int    `json:"dataProviderId"`
	DataConsumer   string `json:"dataConsumer"`
	DataConsumerId int    `json:"dataConsumerId"`
	ShareBeginDate int64  `json:"shareBeginDate"`
	ShareEndDate   int64  `json:"shareEndDate"`
	Interal        uint64 `json:"interal"` //很奇怪的一个词，这里表示获得的积分
	SkuId          int64  `json:"skuId"`
	Timestamp      int64  `json:"timestamp"` //区块链存证时间
	DocType        string `json:"docType"`   //用来区分不同的Struct的，这里的值是ExchangeOrder
}
