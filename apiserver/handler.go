package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func invoke(ctx *gin.Context) {

	request := &Request{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("the received request body is : ", string(requestBytes))

	err = validateRequest(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetPeers := make([]string, 0)
	if request.MetaData.TargetPeers != nil && len(request.MetaData.TargetPeers) != 0 { // high priority
		targetPeers = request.MetaData.TargetPeers
	} else if request.MetaData.TargetOrgs != nil && len(request.MetaData.TargetOrgs) != 0 { // low priority
		for _, orgName := range request.MetaData.TargetOrgs {
			balancer := balancerMgr.OrgBalancers[orgName]
			if balancer == nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "unKnown org " + orgName})
				return
			}
			peer, err := balancer.DoBalance()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			targetPeers = append(targetPeers, peer)
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "targetPeers and targetOrgs both null"})
		return
	}

	channelContext := fabricSdk.ChannelContext(request.MetaData.ChannelName, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	channelClient, err := channel.New(channelContext)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	args := make([][]byte, 0)
	for _, arg := range request.Param.Args {
		args = append(args, []byte(arg))
	}

	result, err := InvokeCC(channelClient,
		channel.Request{
			ChaincodeID: request.MetaData.ChaincodeName,
			Fcn:         request.Param.Function,
			Args:        args,
		},
		targetPeers,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// result.Responses[0].Response.Payload
	status := result.Responses[0].Response.Status
	response := result.Responses[0].Response.Payload

	log.Println("the response is : ", string(response))
	ctx.JSON(http.StatusOK, gin.H{"status": status, "response": string(response)})

}

func query(ctx *gin.Context) {

	request := &Request{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("the received request body is : ", string(requestBytes))

	err = validateRequest(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetPeers := make([]string, 0)
	if request.MetaData.TargetPeers != nil && len(request.MetaData.TargetPeers) != 0 { // high priority
		targetPeers = request.MetaData.TargetPeers
	} else if request.MetaData.TargetOrgs != nil && len(request.MetaData.TargetOrgs) != 0 { // low priority
		for _, orgName := range request.MetaData.TargetOrgs {
			balancer := balancerMgr.OrgBalancers[orgName]
			if balancer == nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "unKnown org " + orgName})
				return
			}
			peer, err := balancer.DoBalance()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			targetPeers = append(targetPeers, peer)
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "targetPeers and targetOrgs both null"})
		return
	}

	channelContext := fabricSdk.ChannelContext(request.MetaData.ChannelName, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	channelClient, err := channel.New(channelContext)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	args := make([][]byte, 0)
	for _, arg := range request.Param.Args {
		args = append(args, []byte(arg))
	}
	result, err := QueryCC(channelClient,
		channel.Request{
			ChaincodeID: request.MetaData.ChaincodeName,
			Fcn:         request.Param.Function,
			Args:        args,
		},
		targetPeers,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status := result.Responses[0].Response.Status
	response := result.Responses[0].Response.Payload

	log.Println("the response is : ", string(response))
	ctx.JSON(http.StatusOK, gin.H{"status": status, "response": string(response)})

}

func saveCustomer(ctx *gin.Context) {

	request := &CustomerInfo{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("the received request body is : ", string(requestBytes))

	channelContext := fabricSdk.ChannelContext("mychannel", fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	channelClient, err := channel.New(channelContext)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	args := make([][]byte, 1)
	args[0] = requestBytes
	targetPeers := make([]string, 1)
	targetPeers[0] = "peer0.org1.example.com"

	result, err := InvokeCC(channelClient,
		channel.Request{
			ChaincodeID: "exchangecc",
			Fcn:         "saveCustomer",
			Args:        args,
		},
		targetPeers,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// result.Responses[0].Response.Payload
	status := result.Responses[0].Response.Status
	response := result.Responses[0].Response.Payload

	data := &SaveDataResponse{
		TransactionHash: string(result.TransactionID),
		//AccountAddress:  strconv.Itoa(request.CustomerId),
	}
	log.Println("the response is : ", string(response))
	ctx.JSON(http.StatusOK, gin.H{"code": status, "message": string(response), "data": data})

}

func saveUploadData(ctx *gin.Context) {

	request := &UploadData{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("the received request body is : ", string(requestBytes))

	channelContext := fabricSdk.ChannelContext("mychannel", fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	channelClient, err := channel.New(channelContext)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	args := make([][]byte, 1)
	args[0] = requestBytes
	targetPeers := make([]string, 1)
	targetPeers[0] = "peer0.org1.example.com"

	result, err := InvokeCC(channelClient,
		channel.Request{
			ChaincodeID: "exchangecc",
			Fcn:         "saveDataUpload",
			Args:        args,
		},
		targetPeers,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// result.Responses[0].Response.Payload
	status := result.Responses[0].Response.Status
	response := result.Responses[0].Response.Payload

	data := &SaveDataResponse{
		TransactionHash: string(result.TransactionID),
		//AccountAddress:  strconv.Itoa(request.CustomerId),
	}
	log.Println("the response is : ", string(response))
	ctx.JSON(http.StatusOK, gin.H{"code": status, "message": string(response), "data": data})

}

func saveExchangeOrder(ctx *gin.Context) {

	request := &ExchangeOrder{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("the received request body is : ", string(requestBytes))

	channelContext := fabricSdk.ChannelContext("mychannel", fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	channelClient, err := channel.New(channelContext)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	args := make([][]byte, 1)
	args[0] = requestBytes
	targetPeers := make([]string, 1)
	targetPeers[0] = "peer0.org1.example.com"

	result, err := InvokeCC(channelClient,
		channel.Request{
			ChaincodeID: "exchangecc",
			Fcn:         "saveExchangeOrder",
			Args:        args,
		},
		targetPeers,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// result.Responses[0].Response.Payload
	status := result.Responses[0].Response.Status
	response := result.Responses[0].Response.Payload

	data := &SaveDataResponse{
		TransactionHash: string(result.TransactionID),
		//AccountAddress:  strconv.Itoa(request.CustomerId),
	}
	log.Println("the response is : ", string(response))
	ctx.JSON(http.StatusOK, gin.H{"code": status, "message": string(response), "data": data})

}

func queryBlock(ctx *gin.Context) {

	channelContext := fabricSdk.ChannelContext("mychannel", fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	client, err := ledger.New(channelContext)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := client.QueryInfo()
	if err != nil {
		log.Fatal(err.Error())
	} else {
		height := result.BCI.Height

		block, err := client.QueryBlock(height - 1)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"blockHash": hex.EncodeToString(block.Header.PreviousHash), "height": height})
	}
}
