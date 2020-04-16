package admin

import (
	"blockchain-sp20/beatchain/chaincode/utils"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var AdminVariable = "admin"

func addProduct(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var err error

	// Access control: Only a Creator can invoke this transaction
	if !utils.AuthenticateCreator(txn) {
		return shim.Error("Caller not a member of Creator Org. Access denied.")
	}

	args := txn.Args
	if len(args) != 2 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 2: {CreatorID, ProductName}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	creatorId := txn.Args[0]
	productName := txn.Args[1]

	// check for valid Creator
	_, err = utils.GetCreatorRecord(stub, txn.Args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	t := time.Now.UnixNano()
	productId := strconv.FormatInt(t, 10) //easiest way for unique identifier

	totalListens := int64(0)
	unRenumeratedListens := int64(0)
	totalMetrics := int64(0)
	unRenumeratedMetrics := int64(0)
	additionalMetrics := int64(0)
	isActive := true

	raw_product := &utils.Product{Id: productId, CreatorId: creatorId, ProductName: productName, TotalListens: totalListens, UnRenumeratedListens: unRenumeratedListens, TotalMetrics: totalMetrics, UnRenumeratedMetrics: unRenumeratedMetrics, AdditionalMetrics: additionalMetrics, IsActive: isActive}
	err = utils.SetProduct(stub, raw_product) //tbd SetProduct
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("Product created cuccessfully with the following attributes :")
	fmt.Println("ProductId : %s", productId)
	fmt.Println("ProductName : %s", productName)
	fmt.Println("CreatorID : %s", creatorId)

	return shim.Success(nil)
}

func deleteProduct(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var err error

	// Access control: Only a Creator can invoke this transaction
	if !utils.AuthenticateCreator(txn) {
		return shim.Error("Caller not a member of Creator Org. Access denied.")
	}

	args := txn.Args
	if len(args) != 2 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 2: {CreatorID, ProductID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	productId := txn.Args[1]

	// check for valid Creator
	creator, err := utils.GetCreatorRecord(stub, txn.Args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	// check for valid Product and verify Creator owns Product
	product, err := utils.GetProduct(stub, txn.Args[1]) //todo getproduct
	if err != nil {
		return shim.Error(err.Error())
	}

	if creator.Id != product.CreatorId {
		err = errors.New(fmt.Sprintf("Creator does not match Product."))
		return shim.Error(err.Error())
	}

	product.IsActive = false

	err = utils.SetProduct(stub, product) //todo setproduct
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Product with id %s has been successfully deleted", productId)

	return shim.Success(nil)
}

func addCustomerRecord(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var err error

	// Access control: Only an admin or app dev can invoke this transaction
	if utils.AuthenticateCreator(txn) || utils.AuthenticateCustomer(txn) {
		return shim.Error("Caller not a member of app dev org or admin org.")
	}

	args := txn.Args
	if len(args) != 2 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 2: {AppDevID, BankAccountId }. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	appDevId := txn.Args[0]
	bankAccountId := txn.Args[1]

	// check for valid AppDev
	_, err = utils.GetAppDevRecord(stub, appDevId)
	if err != nil {
		return shim.Error(err.Error())
	}

	// check for valid bankAccount

	_, err = utils.GetBankAccount(stub, bankAccountId)
	if err != nil {
		return shim.Error(err.Error())
	}

	t := time.Now.UnixNano()
	id := strconv.FormatInt(t, 10) //easiest way for unique identifier

	subscriptionFee := float32(-1) //negative means customer is just created, no subscription as of now

	subscriptionDueDate := time.Time{}  //returns default 'zero' time, the lowest possible one. 

	queuedSong := ""
	previousSong := ""

	raw_customer := &utils.CustomerRecord{Id: id, AppDevId: appDevId, BankAccountId: bankAccountId, SubscriptionFee: subscriptionFee,SubscriptionDueDate: subscriptionDueDate, QueuedSong: queuedSong, PreviousSong: previousSong}
	err = utils.SetCustomerRecord(stub, raw_customer)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Customer successfully created with id: %d",id)
	
	return shim.Success(nil)
}