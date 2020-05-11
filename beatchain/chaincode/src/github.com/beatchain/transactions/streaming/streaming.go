package streaming

import (
	"errors"
	"fmt"
	"time"

	"github.com/beatchain/utils"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var StreamingVariable = "streaming"

/*
	Args = [
		CustomerID,
		AppDevID,
		ProductID,
	]
*/
func RequestSong(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {

	// Access control: Only an Customer Org member can invoke this transaction
	if !txn.TestMode && !(utils.AuthenticateCustomer(txn) || utils.AuthenticateBeatchainAdmin(txn)) {
		return shim.Error("Caller not a member of Customer Org. Access denied.")
	}

	args := txn.Args
	if len(args) != 4 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 4: {CustomerId, AppDevId, CreatorId, ProductID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	customerId := txn.Args[0]
	appDevId := txn.Args[1]
	creatorId := txn.Args[2]
	productId := txn.Args[3]

	customer, err := utils.GetCustomerRecord(stub, customerId)
	if err != nil {
		return shim.Error(err.Error())
	}

	_, err = utils.GetCreatorRecord(stub, creatorId)
	if err != nil {
		return shim.Error(err.Error())
	}

	_, err = utils.GetAppDevRecord(stub, appDevId)
	if err != nil {
		return shim.Error(err.Error())
	}

	product, err := utils.GetProduct(stub, productId)
	if err != nil {
		return shim.Error(err.Error())
	}

	contract, err := utils.GetContract(stub, appDevId, creatorId, productId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if customer.SubscriptionDueDate.After(time.Now()) && customer.AppDevId == appDevId && contract.CreatorId == product.CreatorId {

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

	return shim.Success(nil)
}
