package admin

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/beatchain/utils"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var AdminVariable = "admin"

func AddProduct(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
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

	t := time.Now().UnixNano()
	productId := strconv.FormatInt(t, 10) //easiest way for unique identifier

	totalListens := int64(0)
	unRenumeratedListens := int64(0)
	totalMetrics := int64(0)
	unRenumeratedMetrics := int64(0)
	additionalMetrics := int64(0)
	isActive := true

	rawProduct := &utils.Product{Id: productId, CreatorId: creatorId, ProductName: productName, TotalListens: totalListens, UnRenumeratedListens: unRenumeratedListens, TotalMetrics: totalMetrics, UnRenumeratedMetrics: unRenumeratedMetrics, AdditionalMetrics: additionalMetrics, IsActive: isActive}
	err = utils.SetProduct(stub, rawProduct) //tbd SetProduct
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("Product created cuccessfully with the following attributes :")
	fmt.Println("ProductId : %s", productId)
	fmt.Println("ProductName : %s", productName)
	fmt.Println("CreatorID : %s", creatorId)

	return shim.Success(productId)
}

func DeleteProduct(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
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

func CreateNewBankAccHelper(stub shim.ChaincodeStubInterface) pb.Response {

	var floatBalance float32
	floatBalance = 0.0

	t := time.Now().UnixNano()
	id := strconv.FormatInt(t, 10)

	inUse := true

	rawBankAccount := &utils.BankAccount{Id: id, Balance: floatBalance, InUse: inUse}

	err = utils.SetBankAccount(stub, rawBankAccount)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(id)

}

func AddCustomerRecord(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var err error

	// Access control: Only an admin or app dev can invoke this transaction
	if utils.AuthenticateCreator(txn) || utils.AuthenticateCustomer(txn) {
		return shim.Error("Caller not a member of app dev org or admin org.")
	}

	args := txn.Args
	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {AppDevID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	appDevId := txn.Args[0]

	// check for valid AppDev
	_, err = utils.GetAppDevRecord(stub, appDevId)
	if err != nil {
		return shim.Error(err.Error())
	}

	t := time.Now().UnixNano()
	id := strconv.FormatInt(t, 10) //easiest way for unique identifier

	subscriptionFee := float32(-1) //negative means customer is just created, no subscription as of now

	subscriptionDueDate := time.Time{} //returns default 'zero' time, the lowest possible one.

	queuedSong := ""
	previousSong := ""

	bankAccountId := CreateNewBankAccHelper(stub)

	rawCustomer := &utils.CustomerRecord{Id: id, AppDevId: appDevId, BankAccountId: bankAccountId, SubscriptionFee: subscriptionFee, SubscriptionDueDate: subscriptionDueDate, QueuedSong: queuedSong, PreviousSong: previousSong}
	err = utils.SetCustomerRecord(stub, rawCustomer)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Customer successfully created with id: %s", id)

	return shim.Success(nil)
}

func CreateNewBankAccount(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {

	var err error

	// Access control: Only an admin or app dev can invoke this transaction
	if utils.AuthenticateCreator(txn) || utils.AuthenticateCustomer(txn) || utils.AuthenticateAppDev(txn) {
		return shim.Error("Caller not a member of admin org.")
	}

	args := txn.Args
	if len(args) != 1 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: { Balance }. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	balance := txn.Args[0]
	var floatBalance float32
	var floatBalance64 float64
	if floatBalance64, err = strconv.ParseFloat(balance, 32); err == nil {
		floatBalance = float32(floatBalance64)
		if floatBalance < 0.0 {
			err := errors.New(fmt.Sprintf("Bank Balance must be a positive number \n"))
			return shim.Error(err.Error()) //must be positive
		}
	} else {
		return shim.Error(err.Error())
	}

	t := time.Now().UnixNano()
	id := strconv.FormatInt(t, 10)

	inUse := false

	rawBankAccount := &utils.BankAccount{Id: id, Balance: floatBalance, InUse: inUse}

	err = utils.SetBankAccount(stub, rawBankAccount)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Bank Account successfully created with id: %s and balance: %s ", id, balance)

	return shim.Success(nil)
}

func AddCreatorRecord(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var err error

	// Access control: Only an admin can invoke this transaction
	if utils.AuthenticateCreator(txn) || utils.AuthenticateCustomer(txn) || utils.AuthenticateAppDev(txn) {
		return shim.Error("Caller not a member of admin org.")
	}

	args := txn.Args
	if len(args) != 0 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 0. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	bankAccountId := CreateNewBankAccHelper(stub)

	t := time.Now().UnixNano()
	id := strconv.FormatInt(t, 10)

	rawCreator := &utils.CreatorRecord{Id: id, BankAccountId: bankAccountId}
	err = utils.SetCreatorRecord(stub, rawCreator)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Creator successfully created with id: %s and bank account id: %s and balance = 0.0", id, bankAccountId)

	return shim.Success(id)
}

func AddAppDevRecord(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var account *utils.BankAccount
	var err error

	// Access control: Only an admin can invoke this transaction
	if utils.AuthenticateCreator(txn) || utils.AuthenticateCustomer(txn) || utils.AuthenticateAppDev(txn) {
		return shim.Error("Caller not a member of admin org.")
	}

	args := txn.Args
	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {AdminFeeFrac}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	FeeFrac := txn.Args[1]
	adminFeeFrac64, err := strconv.ParseFloat(FeeFrac, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	adminFeeFrac := float32(adminFeeFrac64)

	if adminFeeFrac < 0.0 || adminFeeFrac > 1.0 {
		err = errors.New(fmt.Sprintf("Admin fee frac must be between 0 and 1"))
		return shim.Error(err.Error())
	}

	t := time.Now().UnixNano()
	id := strconv.FormatInt(t, 10)

	bankAccountId := CreateNewBankAccHelper(stub)

	rawAppDevRecord := &utils.AppDevRecord{Id: id, BankAccountId: bankAccountId, AdminFeeFrac: adminFeeFrac}
	err = utils.SetAppDevRecord(stub, rawAppDevRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Appdev Record successfully created with id: %s, bank account's id %s and admin fee frac : %s", id, bankAccountId, adminFeeFrac)

	return shim.Success(id)
}
