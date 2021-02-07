package main

import "time"

type Request struct {
	MetaData *MetaData `json:"metaData,omitempty" yaml:"metaData,omitempty" binding:"required"`
	Param    *Param    `json:"param,omitempty" yaml:"param,omitempty" binding:"required"`
}

type MetaData struct {
	TargetPeers   []string `json:"targetPeers,omitempty" yaml:"targetPeers,omitempty" `
	TargetOrgs    []string `json:"targetOrgs,omitempty" yaml:"targetOrgs,omitempty" `
	ChannelName   string   `json:"channelName,omitempty" yaml:"channelName,omitempty" valid:"required" binding:"required"`
	ChaincodeName string   `json:"chaincodeName,omitempty" yaml:"chaincodeName,omitempty" valid:"required" binding:"required"`
}

type Param struct {
	Function string   `json:"function,omitempty" yaml:"function,omitempty" valid:"required" binding:"required"`
	Args     []string `json:"args,omitempty" yaml:"args,omitempty" valid:"required" binding:"required"`
}

type CCResponse struct {
	Status  int32  `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Payload string `json:"payload,omitempty"`
}

// ==============================================
type BalancerOrgs struct {
	Organizations map[string]Organiation `json:"organizations,omitempty" yaml:"organizations,omitempty" `
}
type Organiation struct {
	Peers []string `json:"peers,omitempty" yaml:"peers,omitempty" `
}
type SaveDataResponse struct {
	TransactionHash string `json:"transactionHash"`
}

// Block is the summary of block
type Block struct {
	//ChannelName     string    `json:"channel_name"`
	Hash         string `json:"hash"`
	Number       uint64 `json:"number"`
	DataHash     string `json:"data_hash"`
	PreviousHash string `json:"previous_hash"`
	NextHash     string `json:"next_hash"`
	DataSize     int64  `json:"data_size"`
	//KafkaLastOffset int64     `json:"kafka_last_offset"`
	//LastConfigIndex uint64    `json:"last_config_index"`
	TxCount   int       `json:"tx_count"`
	CreatedAt time.Time `json:"created_at"`
}
type Transaction struct {
	//ChannelName      string    `json:"channel_name"`
	ID               string    `json:"id"`
	Type             string    `json:"type"`
	Creator          string    `json:"creator"`
	CreatorMSP       string    `json:"creator_msp"`
	ChaincodeName    string    `json:"chaincode_name"`
	ValidationResult string    `json:"validation_result"`
	BlockNumber      uint64    `json:"block_number"`
	TxNumber         int       `json:"tx_number"`
	CreatedAt        time.Time `json:"created_at"`
}
