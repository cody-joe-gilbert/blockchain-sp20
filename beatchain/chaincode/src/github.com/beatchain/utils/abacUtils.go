package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
)



func GetTxInfo(stub shim.ChaincodeStubInterface, testMode bool) (*Transaction, error) {
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
	var attribute string
	var found bool
	/*
		Construct a more friendly Transaction struct for passing variables
	*/
	txn = new(Transaction)
	txn.CreatorOrg = ""
	txn.CreatorCertIssuer = ""
	txn.LastUniqueId = -1
	txn.TestMode = testMode

	// Fetch the creator org and certificate info
	if !testMode {
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

		// Access Attributes here
		attribute, found, err = cid.GetAttributeValue(stub, "id")
		if found {
			txn.CreatorId = attribute
		} else {
			txn.CreatorId = ""
		}

		// Check for admin rights
		attribute, found, err = cid.GetAttributeValue(stub, "role")
		if found {
			txn.CreatorAdmin = attribute == "admin"
		} else {
			txn.CreatorAdmin = false
		}

	} else {
		// if in test mode, add in dummy values
		txn.CreatorId = "test"
		txn.CreatorOrg = "test"
		txn.CreatorCertIssuer = "test"
		txn.CreatorAdmin = true
	}

	// Fetch the function call info
	txn.CalledFunction, txn.Args = stub.GetFunctionAndParameters()

	return txn, nil
}

/*
The following are helper functions used to authenticate user credentials
 */

func AuthenticateBeatchainAdmin(txn *Transaction) bool {
	return (txn.CreatorOrg == BEATCHAIN_ADMIN_MSP) && (txn.CreatorCertIssuer == BEATCHAIN_ADMIN_CA)
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

