package main

import (
	"errors"
	"flag"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/util/pathvar"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func getFabricConfigProviderFromFile(configPath string) core.ConfigProvider {
	configProvider := config.FromFile(pathvar.Subst(configPath))
	return configProvider
}

func readConfigAndInitSdk() error {
	var err error

	flag.StringVar(&configFabricPath, "configFabric", "./config/config-fabric.yaml", "the path of fabric network config")
	flag.StringVar(&orgName, "orgName", "Org1", "the org name which this apiserver belonged to")
	flag.StringVar(&userName, "uerName", "Admin", "the user name which this apiserver related to")
	flag.StringVar(&balanceType, "balanceType", "robin", "the type of balance")
	flag.Parse()

	log.Println("configFabricPath: ", configFabricPath)
	log.Println("orgName: ", orgName)
	log.Println("userName: ", userName)

	configProvider := getFabricConfigProviderFromFile(configFabricPath)
	fabricSdk, err = fabsdk.New(configProvider)
	if err != nil {
		return err
	}

	err = initBalancerMgr(configFabricPath)
	if err != nil {
		return err
	}

	log.Println("read fabric config success, init fabric sdk success. ")
	log.Println()

	return nil

}

func initBalancerMgr(path string) error {

	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	bOrgs := &BalancerOrgs{}
	err = yaml.Unmarshal(fileData, bOrgs)
	if err != nil {
		return err
	}

	for orgName, org := range bOrgs.Organizations {
		switch balanceType {
		case "robin":
			robinBalancer := &RoundRobinBalancer{
				OrgName: orgName,
				Peers:   org.Peers,
			}
			err := balancerMgr.RegisterBalancer(orgName, robinBalancer)
			if err != nil {
				return err
			}
		case "random":
			randomBalancer := &RandomBalancer{
				OrgName: orgName,
				Peers:   org.Peers,
			}
			err := balancerMgr.RegisterBalancer(orgName, randomBalancer)
			if err != nil {
				return err
			}
		default:
			return errors.New("unKnown balanceType. ")
		}
	}

	return nil
}

func validateRequest(req *Request) error {
	if req == nil {
		return errors.New("request is null. ")
	}

	if (req.MetaData.TargetPeers == nil || len(req.MetaData.TargetPeers) == 0) && len(req.MetaData.TargetOrgs) == 0 {
		return errors.New("targetPeers and targetOrgs both null")
	}

	//if req.Param.Args == nil || len(req.Param.Args) != 1{
	//	return errors.New("args null or len not 1.")
	//}

	return nil
}
