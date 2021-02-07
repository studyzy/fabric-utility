package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"log"

	//	"strconv"
)



var (
	errFormat = "%+v\n"
	logger    = shim.NewLogger("dataShare-cc")
)

// chaincode操作类型
type ChainCode struct {
}

// 初始化方法
func (c *ChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("init complete !")
	return shim.Success(nil)
}

// 外部调用统一入口
func (c *ChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	log.Println()
	logger.Infof("invoke is running %s", function)
	switch function {
	case "regUser"://不关心传入的内容是什么，原封不动的存储到链上
		time,err:=createDataEvidence(stub,args[0],args[1],args[2],args[3],args[4])
		if err!=nil{
			return shim.Error(err.Error())
		}
		return shim.Success(time)
	case "saveData"://以对象的形式存储在链上，方便进行按对象属性的查询
		time,err:=createObjectEvidence(stub,args[0],args[1],args[2],args[3],args[4])
		if err!=nil{
			return shim.Error(err.Error())
		}
		return shim.Success(time)
	case "shareData":
		result,err:=queryByKey(stub, args[0],args[1])
		if err!=nil{
			return shim.Error(err.Error())
		}
		data,err:= json.Marshal(result)
		if err!=nil{
			return shim.Error(err.Error())
		}
		return shim.Success(data)
	case "getAllUsers":
		result,err:=queryByOwner(stub, args[0],args[1])
		if err!=nil{
			return shim.Error(err.Error())
		}
		data,err:= json.Marshal(result)
		if err!=nil{
			return shim.Error(err.Error())
		}
		return shim.Success(data)
	case "getUserById":
		result,err:=queryByOwner(stub, args[0],args[1])
		if err!=nil{
			return shim.Error(err.Error())
		}
		data,err:= json.Marshal(result)
		if err!=nil{
			return shim.Error(err.Error())
		}
		return shim.Success(data)
	case "getDataByKey":
		result,err:=queryByOwner(stub, args[0],args[1])
		if err!=nil{
			return shim.Error(err.Error())
		}
		data,err:= json.Marshal(result)
		if err!=nil{
			return shim.Error(err.Error())
		}
		return shim.Success(data)
	case "getDataShare"://Key:dataKey,userId
		result,err:=queryByOwner(stub, args[0],args[1])
		if err!=nil{
			return shim.Error(err.Error())
		}
		data,err:= json.Marshal(result)
		if err!=nil{
			return shim.Error(err.Error())
		}
		return shim.Success(data)
	default:
		return shim.Error("invalid function:"+function)
	}
}

// chaincode入口
func main() {
	err := shim.Start(new(ChainCode))
	if err != nil {
		logger.Errorf("Error starting chaincode: %s", err)
		return
	}
}