package streaming

import (
	"errors"
	"fmt"
	"time"

	"github.com/beatchain/utils"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


func RequestSong(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {

	// Access control: Only an Customer Org member can invoke this transaction
	if !txn.TestMode && !(utils.AuthenticateCustomer(txn) || utils.AuthenticateBeatchainAdmin(txn)) {
		return shim.Error("Caller not a member of Customer Org. Access denied.")
	}

	if txn.CreatorId == "" {
		return shim.Error("Transaction invoker Customer ID not found in ecert attributes")
	}

	args := txn.Args
	//if len(args) != 4 {
	//	err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 4: {CustomerId, AppDevId, CreatorId, ProductID}. Found %d", len(args)))
	//	return shim.Error(err.Error())
	//}

	if len(args) != 1 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {ProductID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}
	productId := txn.Args[0]

	customer, err := utils.GetCustomerRecord(stub, txn.CreatorId)
	if err != nil {
		return shim.Error(err.Error())
	}

	product, err := utils.GetProduct(stub, productId)
	if err != nil {
		return shim.Error(err.Error())
	}

	appdev, err := utils.GetAppDevRecord(stub, customer.AppDevId)
	if err != nil {
		return shim.Error(err.Error())
	}

	creator, err := utils.GetCreatorRecord(stub, product.CreatorId)
	if err != nil {
		return shim.Error(err.Error())
	}

	contract, err := utils.GetContract(stub, creator.Id, appdev.Id, productId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if customer.SubscriptionDueDate.After(time.Now()) && customer.AppDevId == appdev.Id && contract.CreatorId == product.CreatorId {

		customer.PreviousSong = customer.QueuedSong
		customer.QueuedSong = productId

		err = utils.SetProduct(stub, product)
		if err != nil {
			return shim.Error(err.Error())
		}
	} else {
		err := errors.New(fmt.Sprintf("Invalid combination of parameters or subscription no longer active/valid."))
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("SUCCESS"))
}
