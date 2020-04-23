package main

import (
	"github.com/beatchain/utils"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"time"
)

const (
	layoutISO = "2006-01-02"
)

func ledgerInit(stub shim.ChaincodeStubInterface, txn *utils.Transaction) error {
	/*
	Parses the input variables used to bootstrap the ledger state during first
	initialization

	Args:
		stub: HF shim interface
		txn: parsed Transaction object

	Returns:
		err: Error object. nil if no error occurred.
	 */
	var customerRecord *utils.CustomerRecord
	var appDevRecord *utils.AppDevRecord
	var creatorRecord *utils.CreatorRecord
	var product *utils.Product
	var contract *utils.Contract
	var beatchainAdminBA, customerBA, appDevBA, creatorBA *utils.BankAccount

	// Validate length of arguments
	if len(txn.Args) < 23 {
		return errors.New(fmt.Sprintf("Too few arguments given; given %d arguments", len(txn.Args)))
	}

	fmt.Println("Re-Initializing ledger with input variables: ")

	/*
	Beatchain Admin
	 */
	fmt.Printf("Beatchain Admin BankAccount Initial Balance: %s\n", txn.Args[0])
	beatchainAdminBABalance, err := strconv.ParseFloat(txn.Args[0], 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given beatchainAdminBABalance to float32: %s", txn.Args[0]))
	}

	beatchainAdminBA = &utils.BankAccount{
		Id:      utils.BEATCHAIN_ADMIN_BANK_ACCOUNT_ID,
		Balance: float32(beatchainAdminBABalance),
	}
	err = utils.SetBankAccount(stub, beatchainAdminBA)
	if err != nil {
		return err
	}

	/*
		Test AppDev
	*/
	fmt.Printf("Test AppDev ID: %s\n", txn.Args[1])
	fmt.Printf("Test AppDev BankAccount ID: %s\n", txn.Args[2])
	fmt.Printf("Test AppDev AdminFeeFrac: %s\n", txn.Args[3])
	fmt.Printf("Test AppDev BankAccount Initial Balance: %s\n", txn.Args[4])

	testAppDevId := txn.Args[1]
	testAppDevBAId := txn.Args[2]
	testAppDevAdminFeeFrac, err := strconv.ParseFloat(txn.Args[3], 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testAppDevAdminFeeFrac to float32: %s", txn.Args[3]))
	}
	testAppDevBABalance, err := strconv.ParseFloat(txn.Args[4], 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testAppDevBABalance to float32: %s", txn.Args[4]))
	}

	appDevRecord = &utils.AppDevRecord{
		Id:            testAppDevId,
		BankAccountId: testAppDevBAId,
		AdminFeeFrac:  float32(testAppDevAdminFeeFrac),
	}
	err = utils.SetAppDevRecord(stub, appDevRecord)
	if err != nil {
		return err
	}

	appDevBA = &utils.BankAccount{
		Id:      testAppDevBAId,
		Balance: float32(testAppDevBABalance),
	}
	err = utils.SetBankAccount(stub, appDevBA)
	if err != nil {
		return err
	}

	/*
		Test Customer
	*/
	fmt.Printf("Test Customer ID: %s\n", txn.Args[5])
	fmt.Printf("Test Customer BankAccount ID: %s\n", txn.Args[6])
	fmt.Printf("Test Customer SubscriptionFee: %s\n", txn.Args[7])
	fmt.Printf("Test Customer SubscriptionDueDate: %s\n", txn.Args[8])
	fmt.Printf("Test Customer BankAccount Initial Balance: %s\n", txn.Args[9])

	testCustomerId := txn.Args[5]
	testCustomerBAId := txn.Args[6]
	testCustomerSubscriptionFee, err := strconv.ParseFloat(txn.Args[7], 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testCustomerSubscriptionFee to float32: %s", txn.Args[7]))
	}
	testCustomerSubscriptionDueDate, err := time.Parse(layoutISO, txn.Args[8])
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testCustomerSubscriptionDueDate to date in form YYYY-MM-DD: %s", txn.Args[8]))
	}
	testCustomerBABalance, err := strconv.ParseFloat(txn.Args[9], 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testCustomerBABalance to float32: %s", txn.Args[9]))
	}

	customerBA = &utils.BankAccount{
		Id:      testCustomerBAId,
		Balance: float32(testCustomerBABalance),
	}
	err = utils.SetBankAccount(stub, customerBA)
	if err != nil {
		return err
	}

	customerRecord = &utils.CustomerRecord{
		Id:                  testCustomerId,
		AppDevId:            testAppDevId,
		BankAccountId:       testCustomerBAId,
		SubscriptionFee:     float32(testCustomerSubscriptionFee),
		SubscriptionDueDate: testCustomerSubscriptionDueDate,
		QueuedSong:          "",
		PreviousSong:        "",
	}
	err = utils.SetCustomerRecord(stub, customerRecord)
	if err != nil {
		return err
	}

	/*
		Test Creator
	*/
	fmt.Printf("Test Creator ID: %s\n", txn.Args[10])
	fmt.Printf("Test Creator BankAccount ID: %s\n", txn.Args[11])
	fmt.Printf("Test Creator BankAccount Balance: %s\n", txn.Args[12])

	testCreatorId := txn.Args[10]
	testCreatorBAId := txn.Args[11]
	testCreatorBABalance, err := strconv.ParseFloat(txn.Args[12], 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testCreatorBABalance to float32: %s", txn.Args[12]))
	}

	creatorRecord = &utils.CreatorRecord{
		Id:            testCreatorId,
		BankAccountId: testCreatorBAId,
	}
	err = utils.SetCreatorRecord(stub, creatorRecord)
	if err != nil {
		return err
	}

	creatorBA = &utils.BankAccount{
		Id:      testCreatorBAId,
		Balance: float32(testCreatorBABalance),
	}
	err = utils.SetBankAccount(stub, creatorBA)
	if err != nil {
		return err
	}

	/*
		Test Product
	*/
	fmt.Printf("Test Product ID: %s\n", txn.Args[13])
	fmt.Printf("Test Product Name: %s\n", txn.Args[14])
	fmt.Printf("Test Product Total Listens: %s\n", txn.Args[15])
	fmt.Printf("Test Product Unrenumerated Listens: %s\n", txn.Args[16])
	fmt.Printf("Test Product Total Metrics: %s\n", txn.Args[17])
	fmt.Printf("Test Product Unrenumerated Metrics: %s\n", txn.Args[18])
	fmt.Printf("Test Product Additional Metrics: %s\n", txn.Args[19])
	fmt.Printf("Test Product Status: %s\n", txn.Args[20])

	testProductId := txn.Args[13]
	testProductName := txn.Args[14]
	testProductTotListens, err := strconv.ParseInt(txn.Args[15], 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testProductTotListens to int64: %s", txn.Args[15]))
	}
	testProductUnListens, err := strconv.ParseInt(txn.Args[16], 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testProductUnListens to int64: %s", txn.Args[16]))
	}
	testProductTotMetrics, err := strconv.ParseInt(txn.Args[17], 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testProductTotMetrics to int64: %s", txn.Args[17]))
	}
	testProductUnMetrics, err := strconv.ParseInt(txn.Args[18], 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testProductUnMetrics to int64: %s", txn.Args[18]))
	}
	testProductAddMetrics, err := strconv.ParseInt(txn.Args[19], 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testProductAddMetrics to int64: %s", txn.Args[19]))
	}
	testProductStatus, err := strconv.ParseBool(txn.Args[20])
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testProductStatus to bool: %s", txn.Args[20]))
	}

	product = &utils.Product{
		Id:                   testProductId,
		CreatorId:            testCreatorId,
		ProductName:          testProductName,
		TotalListens:         testProductTotListens,
		UnRenumeratedListens: testProductUnListens,
		TotalMetrics:         testProductTotMetrics,
		UnRenumeratedMetrics: testProductUnMetrics,
		AdditionalMetrics:    testProductAddMetrics,
		IsActive:             testProductStatus,
	}
	err = utils.SetProduct(stub, product)
	if err != nil {
		return err
	}

	/*
		Test Contract
	*/
	fmt.Printf("Test Contract Pay-per-Stream: %s\n", txn.Args[21])
	fmt.Printf("Test Contract Status: %s\n", txn.Args[22])

	testContractPPS, err := strconv.ParseFloat(txn.Args[21], 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given testContractPPS to float32: %s", txn.Args[21]))
	}
	testContractStatus := txn.Args[22]

	contract = &utils.Contract{
		CreatorId:           testCreatorId,
		AppDevId:            testAppDevId,
		ProductId:           testProductId,
		CreatorPayPerStream: float32(testContractPPS),
		Status:              testContractStatus,
	}
	err = utils.SetContract(stub, contract)
	if err != nil {
		return err
	}

	return nil
}
