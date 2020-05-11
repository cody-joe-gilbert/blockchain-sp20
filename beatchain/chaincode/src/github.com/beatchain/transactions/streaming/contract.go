package streaming

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/beatchain/transactions"
	"github.com/beatchain/utils"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var ContractVariable = "contractVariable"

// THESE FUNCTIONS DO NOT CHECK TO SEE IF THE CALLER IS THE ACTUAL ORG MAKING/ACCEPTING/DENYING CONTRACT
// Arun : Can we create a password asset or something like that?
// Julian: This should be something handled through identity management. I'm also not sure how to do it.
func OfferContract(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var err error
	//var creator *utils.CreatorRecord
	//var product *utils.Product

	// Access control: Only an AppDev Org member can invoke this transaction
	if !txn.TestMode && !(utils.AuthenticateAppDev(txn) || utils.AuthenticateBeatchainAdmin(txn)) {
		return shim.Error("Caller not a member of AppDev/Admin Org. Access denied.")
	}

	args := txn.Args
	if len(args) != 4 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 4: {AppDevID, CreatorID, ProductID, CreatorPayPerStream}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	appDevId := txn.Args[0]
	creatorId := txn.Args[1]
	productId := txn.Args[2]
	creatorPayPerStream, err := strconv.ParseFloat(txn.Args[3], 32)
	if err != nil {
		return shim.Error(err.Error())
	}

	// check for valid AppDev
	_, err = utils.GetAppDevRecord(stub, appDevId)
	if err != nil {
		return shim.Error(err.Error())
	}

	// check for valid Creator
	_, err = utils.GetCreatorRecord(stub, creatorId)
	if err != nil {
		return shim.Error(err.Error())
	}

	// check for valid Product
	_, err = utils.GetProduct(stub, productId)
	if err != nil {
		return shim.Error(err.Error())
	}

	//fmt.Println("CreatorId:" + creator.Id + " BankId:" + creator.BankAccountId)
	//fmt.Println("ProductId:" + product.Id + " CreatorId:" + product.CreatorId + " productName:" + product.ProductName)

	//if creator.Id != product.CreatorId {
	//	err = errors.New(fmt.Sprintf("Creator does not match Product. Creator: "  + creator.Id + " Product's Creator: " + product.CreatorId) + " productID: " + product.Id)
	//	return shim.Error(err.Error())
	//}

	raw_contract := &utils.Contract{CreatorId: creatorId, AppDevId: appDevId, ProductId: productId, CreatorPayPerStream: float32(creatorPayPerStream), Status: transactions.REQUESTED}
	err = utils.SetContract(stub, raw_contract)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(transactions.REQUESTED))
}

func AcceptContract(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var err error

	// Access control: Only an Creator Org member can invoke this transaction
	if !txn.TestMode && !(utils.AuthenticateCreator(txn) || utils.AuthenticateBeatchainAdmin(txn)) {
		return shim.Error("Caller not a member of Creator Org. Access denied.")
	}

	args := txn.Args
	if len(args) != 3 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 3: {CreatorID, ProductID, AppDevID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	creatorId := txn.Args[0]
	productId := txn.Args[1]
	appDevId := txn.Args[2]

	contract, err := utils.GetContract(stub, creatorId, appDevId, productId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if contract.Status == transactions.REQUESTED {
		contract.Status = transactions.ACCEPTED
	} else {
		err = errors.New(fmt.Sprintf("Contract has already been finalized."))
		return shim.Error(err.Error())
	}

	err = utils.SetContract(stub, contract)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(transactions.ACCEPTED))
}

func RejectContract(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var err error

	// Access control: Only an Creator Org member can invoke this transaction
	if !txn.TestMode && !(utils.AuthenticateCreator(txn) || utils.AuthenticateBeatchainAdmin(txn)) {
		return shim.Error("Caller not a member of Creator Org. Access denied.")
	}

	args := txn.Args
	if len(args) != 3 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 3: {CreatorID, ProductID, AppDevID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	creatorId := txn.Args[0]
	productId := txn.Args[1]
	appDevId := txn.Args[2]

	contract, err := utils.GetContract(stub, creatorId, appDevId, productId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if contract.Status == transactions.REQUESTED {
		contract.Status = transactions.REJECTED
	} else {
		err = errors.New(fmt.Sprintf("Contract has already been finalized."))
		return shim.Error(err.Error())
	}


	err = utils.SetContract(stub, contract)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(transactions.REJECTED))
}
