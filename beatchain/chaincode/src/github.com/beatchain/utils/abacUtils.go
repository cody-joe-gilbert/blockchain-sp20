package utils

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
)

const CUSTOMER_MSP = "CustomerMSP"
const CUSTOMER_CA = "ca.customerorg.beatchain.com"
const APPDEV_MSP = "AppDevMSP"
const APPDEV_CA = "ca.appdevorg.beatchain.com"
const CREATOR_MSP = "CreatorMSP"
const CREATOR_CA = "ca.creatororg.beatchain.com"

func GetTxInfo(stub shim.ChaincodeStubInterface) (*Transaction, error) {
	/*
	Grabs the transaction info from the calling user

	Args:
		stub: HF shim interface

	Returns:
		txn: Transaction object containing the txn's creator and function call info
		error: Errors raised from accessing attributes. nil if no errors are raised
	 */
	var txn *Transaction
	var certASN1 *pem.Block
	var cert *x509.Certificate
	var err error
	var attribute string
	var found bool
	/*
		Construct a more friendly Transaction struct for passing variables
	*/
	txn = new(Transaction)
	txn.CreatorOrg = ""
	txn.CreatorCertIssuer = ""
	txn.TestMode = false

	// Fetch the creator org and certificate info
	creator, err := stub.GetCreator()
	if err != nil {
		_ = fmt.Errorf("Error getting transaction creator: %s\n", err.Error())
		return txn, err
	}
	creatorSerializedId := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creator, creatorSerializedId)
	if err != nil {
		fmt.Printf("Error unmarshalling creator identity: %s\n", err.Error())
		return txn, err
	}
	if len(creatorSerializedId.IdBytes) == 0 {
		return txn, errors.New("empty certificate")
	}
	certASN1, _ = pem.Decode(creatorSerializedId.IdBytes)
	cert, err = x509.ParseCertificate(certASN1.Bytes)
	if err != nil {
		return txn, err
	}
	txn.CreatorOrg = creatorSerializedId.Mspid
	txn.CreatorCertIssuer = cert.Issuer.CommonName

	// Fetch the function call info
	txn.CalledFunction, txn.Args = stub.GetFunctionAndParameters()

	// Access Attributes here
	attribute, found, err = cid.GetAttributeValue(stub, "id")
	if found {
		txn.CreatorId = attribute
	} else {
		txn.CreatorId = ""
	}

	return txn, nil
}

func AuthenticateCustomer(txn *Transaction) bool {
	return (txn.CreatorOrg == CUSTOMER_MSP) && (txn.CreatorCertIssuer == CUSTOMER_CA)
}

func AuthenticateAppDev(txn *Transaction) bool {
	return (txn.CreatorOrg == APPDEV_MSP) && (txn.CreatorCertIssuer == APPDEV_CA)
}

func AuthenticateCreator(txn *Transaction) bool {
	return (txn.CreatorOrg == CREATOR_MSP) && (txn.CreatorCertIssuer == CREATOR_CA)
}

