package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//创建存证
func createDataEvidence(stub shim.ChaincodeStubInterface, category,owner,dataKey,dataValue,reference string) ([]byte,error) {
	t,err:=stub.GetTxTimestamp()
	if err != nil{
		return nil,err
	}
	evidence:=&DataEvidence{
		Owner:     owner,
		Category:  category,
		DataKey:   dataKey,
		DataValue: dataValue,
		Reference: reference,
		Timestamp: t.Seconds,
	}
	key:=fmt.Sprintf("D_%s_%s",category,dataKey)

	bytes, err := json.Marshal(evidence)
	if err != nil{
		return nil, err
	}
	err = stub.PutState(key,bytes)
	if err != nil{
		return nil, err
	}


	return []byte(fmt.Sprintf("%d",t.Seconds)),nil
}
func createObjectEvidence(stub shim.ChaincodeStubInterface, category,owner,dataKey,objectJson,reference string) ([]byte,error) {
	t,err:=stub.GetTxTimestamp()
	if err != nil{
		return nil,err
	}
	obj:=make(map[string]interface{})
	err=json.Unmarshal([]byte(objectJson),&obj)
	if err != nil{
		return nil,err
	}
	evidence:=&ObjectEvidence{
		Owner:     owner,
		Category:  category,
		ObjectKey:   dataKey,
		Object: obj,
		Reference: reference,
		Timestamp: t.Seconds,
	}
	key:=fmt.Sprintf("O_%s_%s",category,dataKey)

	bytes, err := json.Marshal(evidence)
	if err != nil{
		return nil, err
	}
	err = stub.PutState(key,bytes)
	if err != nil{
		return nil, err
	}


	return []byte(fmt.Sprintf("%d",t.Seconds)),nil
}


func queryByKey(stub shim.ChaincodeStubInterface, category,dataKey string) (*DataEvidence,error) {
	key:=fmt.Sprintf("%s_%s",category,dataKey)
	bytes, err:= stub.GetState(key)
	if err != nil{
		return nil, err
	}
	evidence:=&DataEvidence{}
	err=json.Unmarshal(bytes,evidence)
	if err != nil{
		return nil, err
	}
	return evidence,nil
}

func queryByOwner(stub shim.ChaincodeStubInterface, category,owner string) ([]*DataEvidence,error) {
	queryStr:=fmt.Sprintf("{\"selector\":{\"owner\":\"%s\"}}",owner)
	resultsIterator, err:= stub.GetQueryResult(queryStr)
	defer resultsIterator.Close()
	if err != nil {
		return nil, err
	}
	result:=[]*DataEvidence{}
		for resultsIterator.HasNext() {
		queryResponse,err:= resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		evi:=&DataEvidence{}
		json.Unmarshal(queryResponse.Value,evi)
		result=append(result,evi)
	}
	return result, nil
}
