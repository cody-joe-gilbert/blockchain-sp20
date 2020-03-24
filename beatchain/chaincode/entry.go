package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	/*
		Import transaction packages here
	*/
	"blockchain-sp20/beatchain/chaincode/transactions"
	"blockchain-sp20/beatchain/chaincode/transactions/admin"
	"blockchain-sp20/beatchain/chaincode/transactions/banking"
	"blockchain-sp20/beatchain/chaincode/transactions/streaming"

)

// BeatchainChaincode implementation
type BeatchainChaincode struct {
	testMode bool
}

// Initialization template
func (t *BeatchainChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Initializing Beatchain")
	_, args := stub.GetFunctionAndParameters()
	var err error
	/*
		Typechecking and initialization here
	*/
	return shim.Success(nil)
}

// Invocation template
func (t *TradeWorkflowChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("BeatchainChaincode Invoke")

	creator, err := stub.GetCreator()
	if err != nil {
		fmt.Errorf("Error getting transaction creator: %s\n", err.Error())
		return shim.Error(err.Error())
	}
	creatorOrg := ""
	creatorCertIssuer := ""
	if !t.testMode {
		creatorOrg, creatorCertIssuer, err = getTxCreatorInfo(creator)
		if err != nil {
			fmt.Errorf("Error extracting creator identity info: %s\n", err.Error())
			return shim.Error(err.Error())
		}
		fmt.Printf("TradeWorkflow Invoke by '%s', '%s'\n", creatorOrg, creatorCertIssuer)
	}

	fmt.Println("This is a call to the %v module and file.", transactions.TransactionsVariable)
	fmt.Println("This is a call to the %v module and file.", admin.AdminVariable)
	fmt.Println("This is a call to the %v module and file.", banking.BankingVariable)
	fmt.Println("This is a call to the %v module and file.", streaming.StreamingVariable)

	/*
		Here we'll dispatch invocation to separate function modules
	*/
	function, args := stub.GetFunctionAndParameters()
	if function == "" {
		// Importer requests a trade
		return t.FUNCTION(stub, creatorOrg, creatorCertIssuer, args)
	}
	return shim.Error("Invalid invoke function name")
}
