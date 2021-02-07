package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	orgName          string
	userName         string
	configFabricPath string
	fabricSdk        *fabsdk.FabricSDK
	balanceType      string
)

var balancerMgr = &BalancerMgr{
	OrgBalancers: make(map[string]Balancer, 0),
}

func main() {

	err := readConfigAndInitSdk()
	if err != nil {
		panic(err)
	}
	defer fabricSdk.Close()

	//ticker := time.NewTicker(time.Second * 2)
	//channelContext := fabricSdk.ChannelContext("mychannel", fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	//client, _ := ledger.New(channelContext)
	//i := 0
	//go func() {
	//	for { //循环
	//		<-ticker.C
	//		i++
	//		//读取最高区块高度
	//
	//		result,err:= client.QueryInfo()
	//		if err!=nil{
	//			log.Fatal(err.Error())
	//		}else {
	//			height := result.BCI.Height
	//
	//			block,err:=client.QueryBlock(height-1)
	//			if err!=nil{
	//				log.Print(err)
	//			}else{
	//				log.Printf("%d, Height:%d, parentHash:%x",i, height,block.Header.PreviousHash)
	//			}
	//		}
	//	}
	//}()

	router := gin.Default()
	api := router.Group("/api")
	{

		api.POST("/customer", saveCustomer)
		api.POST("/uploadData", saveUploadData)
		api.POST("/exchangeOrder", saveExchangeOrder)
		api.GET("/block", queryBlock)
		api.GET("/query", query)
		api.POST("/invoke", invoke)
	}

	router.Run(":1001")
}
