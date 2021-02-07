package main

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric/protos/common"
	mb "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
	putils "github.com/hyperledger/fabric/protos/utils"
	"github.com/studyzy/fabric-utility/apiserver/cache"
)

// cache identity
var identities cache.Store

func init() {
	identities = cache.NewLRU(100)
}
func convertToBlock(block *cb.Block) (*Block, error) {
	blk := &Block{
		//ChannelName:     channelName,
		Hash:         hex.EncodeToString(block.Header.Hash()),
		Number:       block.Header.Number,
		DataHash:     hex.EncodeToString(block.Header.DataHash),
		PreviousHash: hex.EncodeToString(block.Header.PreviousHash),
		DataSize:     int64(len(putils.MarshalOrPanic(block))),
		//KafkaLastOffset: getKafkaOffsetMetadataFromBlock(block),
		//LastConfigIndex: putils.GetLastConfigIndexFromBlockOrPanic(block),
		TxCount: len(block.Data.Data),
	}

	txFlags := block.Metadata.Metadata[cb.BlockMetadataIndex_TRANSACTIONS_FILTER]
	lastEnv := unmarshalEnvelopeOrPanic(block.Data.Data[blk.TxCount-1])

	tx, err := convertEnvelopeToTx(int32(txFlags[blk.TxCount-1]), lastEnv)
	if err != nil {
		return nil, err
	}

	blk.CreatedAt = tx.CreatedAt

	return blk, nil
}
func unmarshalEnvelopeOrPanic(b []byte) *cb.Envelope {
	env := &cb.Envelope{}
	if err := proto.Unmarshal(b, env); err != nil {
		panic(fmt.Errorf("Error getting envelope(%s)", err))
	}

	return env
}

func convertEnvelopeToTx(txFlag int32, env *cb.Envelope) (*Transaction, error) {
	payload, err := putils.GetPayload(env)
	if err != nil {
		return nil, fmt.Errorf("Unexpected error from payload: %v", err)
	}

	chdr, err := putils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, fmt.Errorf("Unexpected error from unmarshal channel header: %v", err)
	}
	shdr, err := putils.GetSignatureHeader(payload.Header.SignatureHeader)
	if err != nil {
		return nil, fmt.Errorf("Unexpected error from unmarshal signature header: %v", err)
	}
	identity, err := getIdentity(shdr.Creator)
	if err != nil {
		return nil, err
	}
	tx := &Transaction{
		ID:               chdr.TxId,
		Type:             cb.HeaderType_name[chdr.Type],
		CreatorMSP:       identity.mspID,
		ValidationResult: pb.TxValidationCode_name[int32(txFlag)],
		CreatedAt:        time.Unix(chdr.Timestamp.Seconds, int64(chdr.Timestamp.Nanos)),
	}

	// Ecert
	//if identity.cert != nil {
	//	tx.Creator = identity.cert.Subject.CommonName
	//}

	hdrExt, err := putils.GetChaincodeHeaderExtension(payload.Header)
	if err != nil {
		return nil, fmt.Errorf("GetChaincodeHeaderExtension failed: %v", err)
	}

	if hdrExt.ChaincodeId != nil {
		tx.ChaincodeName = hdrExt.ChaincodeId.Name
	}

	return tx, nil
}

type cachedIdentity struct {
	mspID string
	cert  *x509.Certificate
}

func getIdentity(serilizedIdentity []byte) (*cachedIdentity, error) {
	var err error
	newIdentity := func() interface{} {
		sid := &mb.SerializedIdentity{}
		err = proto.Unmarshal(serilizedIdentity, sid)
		if err != nil {
			return nil
		}

		var cert *x509.Certificate
		cert, err = decodeX509Pem(sid.IdBytes)
		if err != nil {
			return nil
		}

		return &cachedIdentity{
			mspID: sid.Mspid,
			cert:  cert,
		}
	}

	identity, _ := identities.SetNotExist(string(serilizedIdentity), newIdentity)
	if err != nil {
		return nil, err
	}

	return identity.(*cachedIdentity), nil

}

func decodeX509Pem(certPem []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certPem)
	if block == nil {
		return nil, fmt.Errorf("bad cert")
	}

	return x509.ParseCertificate(block.Bytes)
}
