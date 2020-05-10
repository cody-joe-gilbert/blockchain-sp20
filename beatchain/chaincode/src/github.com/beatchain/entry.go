package main

import (
	"fmt"

	"github.com/beatchain/transactions/streaming"
	"github.com/beatchain/transactions/admin"
	"github.com/beatchain/transactions/banking"
	"github.com/beatchain/utils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// BeatchainChaincode implementation
type BeatchainChaincode struct {
	testMode bool
}

// Initialization template
func (t *BeatchainChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	var txn *utils.Transaction
	var err error

	fmt.Println("Initializing Beatchain chaincode")

	// Get the transaction details
	txn, err = utils.GetTxInfo(stub, t.testMode)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(txn.Args) == 0 {
		// Using existing ledger
		fmt.Println("Initializing with existing ledger")
		return shim.Success(nil)
	}

	// New variables given; initialize ledger
	err = ledgerInit(stub, txn)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Invocation template
func (t *BeatchainChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	var txn *utils.Transaction
	var err error

	fmt.Println("BeatchainChaincode Invoke")

	// Get the transaction details

	txn, err = utils.GetTxInfo(stub, t.testMode)
	if err != nil {
		return shim.Error(err.Error())
	}

	/*
		Here we'll dispatch invocation to separate function modules
	*/
	switch fnct := txn.CalledFunction; fnct {

	case "ListBankAccounts":
		return admin.ListBankAccounts(stub, txn)
	case "ListAllCustomers":
		return admin.ListAllCustomers(stub, txn)
	case "ListAppCustomers":
		return admin.ListAppCustomers(stub, txn)
	case "AddProduct":
		return admin.AddProduct(stub, txn)
	case "DeleteProduct":
		return admin.DeleteProduct(stub, txn)
	case "AddCustomerRecord":
		return admin.AddCustomerRecord(stub, txn)
	case "AddAppDevRecord":
		return admin.AddAppDevRecord(stub, txn)
	case "AddCreatorRecord":
		return admin.AddCreatorRecord(stub, txn)
	case "RenewSubscription":
		return banking.RenewSubscription(stub, txn)
	case "CollectPayment":
		return banking.CollectPayment(stub, txn)
	case "TransferFunds":
		return banking.TransferFunds(stub, txn)
	case "OfferContract":
		return streaming.OfferContract(stub, txn)
	case "AcceptContract":
		return streaming.AcceptContract(stub, txn)
	case "RejectContract":
		return streaming.RejectContract(stub, txn)
	case "RequestSong":
		return streaming.RequestSong(stub, txn)
	default:
		return shim.Error("Invalid invoke function name")
	}

}

func main() {
	/*
		Bootstraps the Beatchain chaincode
	*/
	bcc := new(BeatchainChaincode)
	bcc.testMode = true
	err := shim.Start(bcc)
	if err != nil {
		fmt.Printf("Error starting Trade Workflow chaincode: %s", err)
	}
}
