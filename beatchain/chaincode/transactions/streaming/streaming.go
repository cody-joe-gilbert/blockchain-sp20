package streaming

import (
	"blockchain-sp20/beatchain/chaincode/utils"
	"errors"
	"fmt"

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
func requestSong(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {

	// Access control: Only an Customer Org member can invoke this transaction
	if !utils.AuthenticateCustomer(txn) {
		return shim.Error("Caller not a member of Customer Org. Access denied.")
	}

	args := txn.Args
	if len(args) != 3 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 3: {CustomerID, AppDevID, ProductID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	customerID := txn.Args[0]
	appDevId := txn.Args[1]
	productId := txn.Args[2]

	customerKey, err := utils.GetCustomerRecordKey(stub, customerID)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println(appDevId)
	fmt.Println(productId)
	fmt.Println(customerKey)

	return shim.Success(nil)
}
