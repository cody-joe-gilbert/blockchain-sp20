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
func OfferContract(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	/*
	Offers a contract with a given payment per stream in $USD to a Creator for the rights
	to stream content.

	Args:
		AppDevID (string): ID of the AppDev Submitting the offer
			TODO: Replace with txn.CreatorId
		CreatorID (string): ID of the Creator to which the contract is offered
		ProductID (string): ID of the Product under consideration of the contract
			Note: Each Product has a separate contract in this draft
		CreatorPayPerStream (float32): Payment in $USD per stream of the product
	*/

	var creator *utils.CreatorRecord
	var product *utils.Product
	var err error

	// Access control: Only an AppDev Org member can invoke this transaction
	//if !utils.AuthenticateAppDev(txn) {
	//	return shim.Error("Caller not a member of AppDev Org. Access denied.")
	//}

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
	_, err = utils.GetAppDevRecord(stub, txn.Args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	// check for valid Creator
	creator, err = utils.GetCreatorRecord(stub, txn.Args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	// check for valid Product and verify Creator owns Product
	product, err = utils.GetProduct(stub, txn.Args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	if creator.Id != product.CreatorId {
		err = errors.New(fmt.Sprintf("Creator does not match Product. Creator: "  + creator.Id + " Product's Creator: " + product.CreatorId) + " productID: " + product.Id)
		return shim.Error(err.Error())
	}

	raw_contract := &utils.Contract{CreatorId: creatorId, AppDevId: appDevId, ProductId: productId, CreatorPayPerStream: float32(creatorPayPerStream), Status: transactions.REQUESTED}
	err = utils.SetContract(stub, raw_contract)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("SUCCESS"))
}

func AcceptContract(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	/*
		Accepts an offered contract for a given payment per stream in $USD by a Creator for the rights
		to stream content.

		Args:
			CreatorID (string): ID of the Creator to which the contract is offered
				TODO: Replace with txn.CreatorId
			ProductID (string): ID of the Product under consideration of the contract
				Note: Each Product has a separate contract in this draft
			AppDevID (string): ID of the AppDev Submitting the offer
	*/

	var err error

	// Access control: Only an Creator Org member can invoke this transaction
	if !utils.AuthenticateCreator(txn) {
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

	contract.Status = transactions.ACCEPTED

	err = utils.SetContract(stub, contract)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("SUCCESS"))
}

func RejectContract(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	/*
		Rejects an offered contract by a Creator for the rights to stream content.

		Args:
			CreatorID (string): ID of the Creator to which the contract is offered
				TODO: Replace with txn.CreatorId
			ProductID (string): ID of the Product under consideration of the contract
				Note: Each Product has a separate contract in this draft
			AppDevID (string): ID of the AppDev Submitting the offer
	*/
	var err error

	// Access control: Only an Creator Org member can invoke this transaction
	if !utils.AuthenticateCreator(txn) {
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

	contract.Status = transactions.REJECTED

	err = utils.SetContract(stub, contract)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("SUCCESS"))
}
