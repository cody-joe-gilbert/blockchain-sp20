package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	/*
		Import transaction packages here
	*/
	"blockchain-sp20/beatchain/chaincode/transactions"
	"blockchain-sp20/beatchain/chaincode/transactions/admin"
	"blockchain-sp20/beatchain/chaincode/transactions/banking"
	"blockchain-sp20/beatchain/chaincode/transactions/streaming"

	"blockchain-sp20/beatchain/chaincode/utils"

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
	txn, err = utils.GetTxInfo(stub)
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
	//var txn *utils.Transaction
	//var err error

	fmt.Println("BeatchainChaincode Invoke")

	// Get the transaction details
/*	txn, err = utils.GetTxInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}*/

	fmt.Println("This is a call to the %v module and file.", transactions.TransactionsVariable)
	fmt.Println("This is a call to the %v module and file.", admin.AdminVariable)
	fmt.Println("This is a call to the %v module and file.", banking.BankingVariable)
	fmt.Println("This is a call to the %v module and file.", streaming.StreamingVariable)

	/*
		Here we'll dispatch invocation to separate function modules
	*/

	/*if txn.CalledFunction == "" {
		// Importer requests a trade
		return t.FUNCTION(stub, txn)
	}*/
	return shim.Error("Invalid invoke function name")
}

func main() {
	/*
	Bootstraps the Beatchain chaincode
	 */
	bcc := new(BeatchainChaincode)
	bcc.testMode = false
	err := shim.Start(bcc)
	if err != nil {
		fmt.Printf("Error starting Trade Workflow chaincode: %s", err)
	}
}