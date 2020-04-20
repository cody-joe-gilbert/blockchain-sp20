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
	var beatchainAdminBA, customerBA, appDevBA *utils.BankAccount

	// Validate length of arguments
	if len(txn.Args) < 10 {
		return errors.New(fmt.Sprintf("Too few arguments given; given %d arguments", len(txn.Args)))
	}

	fmt.Println("Re-Initializing ledger with input variables: ")

	fmt.Printf("Beatchain Admin BankAccount Initial Balance: %s\n", txn.Args[0])

	fmt.Printf("Test AppDev ID: %s\n", txn.Args[1])
	fmt.Printf("Test AppDev BankAccount ID: %s\n", txn.Args[2])
	fmt.Printf("Test AppDev AdminFeeFrac: %s\n", txn.Args[3])
	fmt.Printf("Test AppDev BankAccount Initial Balance: %s\n", txn.Args[4])

	fmt.Printf("Test Customer ID: %s\n", txn.Args[5])
	fmt.Printf("Test Customer BankAccount ID: %s\n", txn.Args[6])
	fmt.Printf("Test Customer SubscriptionFee: %s\n", txn.Args[7])
	fmt.Printf("Test Customer SubscriptionDueDate: %s\n", txn.Args[8])
	fmt.Printf("Test Customer BankAccount Initial Balance: %s\n", txn.Args[9])

	// Parse out and validate all of the above input
	beatchainAdminBABalance, err := strconv.ParseFloat(txn.Args[0], 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse given beatchainAdminBABalance to float32: %s", txn.Args[0]))
	}

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

	// Setup ledger objects
	beatchainAdminBA = &utils.BankAccount{
		Id:      utils.BEATCHAIN_ADMIN_BANK_ACCOUNT_ID,
		Balance: float32(beatchainAdminBABalance),
	}

	appDevBA = &utils.BankAccount{
		Id:      testAppDevBAId,
		Balance: float32(testAppDevBABalance),
	}

	customerBA = &utils.BankAccount{
		Id:      testCustomerBAId,
		Balance: float32(testCustomerBABalance),
	}

	appDevRecord = &utils.AppDevRecord{
		Id:            testAppDevId,
		BankAccountId: testAppDevBAId,
		AdminFeeFrac:  float32(testAppDevAdminFeeFrac),
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

	// Set objects in the ledger
	err = utils.SetBankAccount(stub, beatchainAdminBA)
	if err != nil {
		return err
	}
	err = utils.SetBankAccount(stub, appDevBA)
	if err != nil {
		return err
	}
	err = utils.SetBankAccount(stub, customerBA)
	if err != nil {
		return err
	}
	err = utils.SetCustomerRecord(stub, customerRecord)
	if err != nil {
		return err
	}
	err = utils.SetAppDevRecord(stub, appDevRecord)
	if err != nil {
		return err
	}

	return nil
}
