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


func AddProduct(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var id string
	var err error

	// Access control: Only a Creator can invoke this transaction
	if !txn.TestMode && !utils.AuthenticateCreator(txn) {
		return shim.Error("Caller not a member of Creator Org. Access denied.")
	}

	if txn.CreatorId == "" {
		return shim.Error("Transaction invoker Creator ID not found in ecert attributes")
	}

	if len(txn.Args) != 1 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {ProductName}. Found %d", len(txn.Args)))
		return shim.Error(err.Error())
	}

	// check for valid Creator
	_, err = utils.GetCreatorRecord(stub, txn.CreatorId)
	if !txn.TestMode && err != nil {
		return shim.Error(err.Error())
	}
	// Get a unique key
	id, err = utils.GetUniqueId(stub, txn)
	if err != nil {
		return shim.Error(err.Error())
	}

	rawProduct := &utils.Product{Id: id,
		CreatorId: txn.CreatorId,
		ProductName: txn.Args[0],
		TotalListens: int64(0),
		UnRenumeratedListens: int64(0),
		TotalMetrics: int64(0),
		UnRenumeratedMetrics: int64(0),
		AdditionalMetrics: int64(0),
		IsActive: true}

	err = utils.SetProduct(stub, rawProduct) //tbd SetProduct
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("Product created successfully with the following attributes :")
	fmt.Printf("ProductId : %s", id)
	fmt.Printf("ProductName : %s", txn.Args[0])
	fmt.Printf("CreatorID : %s", txn.CreatorId)

	return shim.Success([]byte(id))
}

func DeleteProduct(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var err error

	// Access control: Only a Creator can invoke this transaction
	if !txn.TestMode && !utils.AuthenticateCreator(txn) {
		return shim.Error("Caller not a member of Creator Org. Access denied.")
	}

	if txn.CreatorId == "" {
		return shim.Error("Transaction invoker Creator ID not found in ecert attributes")
	}

	if len(txn.Args) != 1 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {ProductID}. Found %d", len(txn.Args)))
		return shim.Error(err.Error())
	}

	productId := txn.Args[0]

	// check for valid Creator
	creator, err := utils.GetCreatorRecord(stub, txn.CreatorId)
	if err != nil {
		return shim.Error(err.Error())
	}

	// check for valid Product and verify Creator owns Product
	product, err := utils.GetProduct(stub, txn.Args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	if creator.Id != product.CreatorId {
		err = errors.New(fmt.Sprintf("Creator ID %s does not match Product's creator id %s", creator.Id, product.CreatorId))
		return shim.Error(err.Error())
	}

	product.IsActive = false

	err = utils.SetProduct(stub, product)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Product with id %s has been successfully deleted", productId)

	return shim.Success([]byte("SUCCESS"))
}

func createNewBankAccHelper(stub shim.ChaincodeStubInterface, txn *utils.Transaction) (string, error) {
	var id string
	var floatBalance float32
	var err error
	floatBalance = 0.0

	// Get a unique key
	id, err = utils.GetUniqueId(stub, txn)
	if err != nil {
		return id, err
	}

	rawBankAccount := &utils.BankAccount{Id: id, Balance: floatBalance, InUse: true}

	err = utils.SetBankAccount(stub, rawBankAccount)
	if err != nil {
		return "", err
	}

	return id, nil
}

func AddCustomerRecord(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var id string
	var err error

	// Access control: Only appdev org can invoke this transaction
	if !txn.TestMode && !utils.AuthenticateAppDev(txn) {
		return shim.Error("Caller not a member of appdev org.")
	}

	if txn.CreatorId == "" {
		return shim.Error("Transaction invoker AppDev ID not found in ecert attributes")
	}

	if len(txn.Args) != 1 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {subscriptionFee}. Found %d", len(txn.Args)))
		return shim.Error(err.Error())
	}

	subscriptionFee, err := strconv.ParseFloat(txn.Args[0], 32)
	if err != nil {
		err = errors.New(fmt.Sprintf("Cannot parse given subscriptionFee to float32: %s", txn.Args[0]))
		return shim.Error(err.Error())
	}

	// check for valid AppDev
	_, err = utils.GetAppDevRecord(stub, txn.CreatorId)
	if !txn.TestMode && err != nil {
		fmt.Printf("Cannot find AppDevRecord with ID %s", txn.CreatorId)
		return shim.Error(err.Error())
	}

	// Get a unique key
	id, err = utils.GetUniqueId(stub, txn)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Sub due 1 month from creation date
	t := time.Now()
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	subscriptionDueDate := date.Add(time.Hour * 24 * 30)

	bankAccountId, err := createNewBankAccHelper(stub, txn)
	if err != nil {
		return shim.Error(err.Error())
	}

	rawCustomer := &utils.CustomerRecord{
		Id: id,
		AppDevId: txn.CreatorId,
		BankAccountId: bankAccountId,
		SubscriptionFee: float32(subscriptionFee),
		SubscriptionDueDate: subscriptionDueDate,
		QueuedSong: "",
		PreviousSong: ""}

	err = utils.SetCustomerRecord(stub, rawCustomer)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Customer successfully created with id: %s", id)

	return shim.Success([]byte(id))
}

func CreateNewBankAccount(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var id string
	var err error

	// Org-based Access control not needed as all orgs access this function

	// Check for admin rights
	if !txn.CreatorAdmin {
		err := errors.New(fmt.Sprint("access denied: function requires admin privileges"))
		return shim.Error(err.Error())
	}

	if len(txn.Args) != 1 {
		err := errors.New(fmt.Sprintf("Expecting 0 arguments. Found %d", len(txn.Args)))
		return shim.Error(err.Error())
	}


	// Get a unique key
	id, err = utils.GetUniqueId(stub, txn)
	if err != nil {
		return shim.Error(err.Error())
	}

	// All BAs initialized to $0.00 to prevent money creation via account creation
	rawBankAccount := &utils.BankAccount{Id: id, Balance: 0.00, InUse: false}

	err = utils.SetBankAccount(stub, rawBankAccount)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Bank Account successfully created with id: %s", id)

	return shim.Success([]byte(id))
}

func AddCreatorRecord(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var id string
	var err error

	// Check for admin rights
	if !txn.CreatorAdmin {
		err := errors.New(fmt.Sprint("access denied: function requires admin privileges"))
		return shim.Error(err.Error())
	}

	// Access control: Only a creator org can invoke this transaction
	if !txn.TestMode && !utils.AuthenticateCreator(txn) {
		return shim.Error("Caller not a member of Creator org.")
	}

	if len(txn.Args) != 0 {
		err := errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 0. Found %d", len(txn.Args)))
		return shim.Error(err.Error())
	}

	bankAccountId, err := createNewBankAccHelper(stub, txn)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Get a unique key
	id, err = utils.GetUniqueId(stub, txn)
	if err != nil {
		return shim.Error(err.Error())
	}

	rawCreator := &utils.CreatorRecord{Id: id, BankAccountId: bankAccountId}
	err = utils.SetCreatorRecord(stub, rawCreator)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Creator successfully created with id: %s and bank account id: %s and balance = 0.0", id, bankAccountId)

	return shim.Success([]byte(id))
}

func AddAppDevRecord(stub shim.ChaincodeStubInterface, txn *utils.Transaction) pb.Response {
	var id string
	var err error

	if !txn.CreatorAdmin {
		err := errors.New(fmt.Sprint("access denied: function requires admin privileges"))
		return shim.Error(err.Error())
	}

	// Access control: Only an admin can invoke this transaction
	if !txn.TestMode && !utils.AuthenticateAppDev(txn) {
		return shim.Error("Caller not a member of appdev org.")
	}

	if len(txn.Args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {AdminFeeFrac}. Found %d", len(txn.Args)))
		return shim.Error(err.Error())
	}

	adminFeeFrac64, err := strconv.ParseFloat(txn.Args[0], 32)
	if err != nil {
		err = errors.New(fmt.Sprintf("Cannot parse given adminFeeFrac to float32: %s", txn.Args[0]))
		return shim.Error(err.Error())
	}
	adminFeeFrac := float32(adminFeeFrac64)

	if adminFeeFrac < 0.0 || adminFeeFrac > 1.0 {
		err = errors.New(fmt.Sprintf("Admin fee frac must be between 0 and 1"))
		return shim.Error(err.Error())
	}


	// Get a unique key


	bankAccountId, err := createNewBankAccHelper(stub, txn)
	if err != nil {
		return shim.Error(err.Error())
	}

	id, err = utils.GetUniqueId(stub, txn)
	if err != nil {
		return shim.Error(err.Error() + "; id:" + id)
	}

	rawAppDevRecord := &utils.AppDevRecord{Id: id, BankAccountId: bankAccountId, AdminFeeFrac: adminFeeFrac}
	err = utils.SetAppDevRecord(stub, rawAppDevRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Appdev Record successfully created with id: %s, bank accounts id %s and admin fee frac : %.2f", id, bankAccountId, adminFeeFrac)

	return shim.Success([]byte(id))
}