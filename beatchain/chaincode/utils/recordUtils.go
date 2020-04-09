package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func GetCustomerRecord(stub shim.ChaincodeStubInterface, customerId string) (*CustomerRecord, error) {
	/*
		Fetches a CustomerRecord object from off the ledger

		Args:
			stub: HF shim interface
			customerId: Primary Key of the Customer

		Returns:
			customerRecord: CustomerRecord struct obj for the requested record
			err: Error object. nil if no error occurred.

	*/
	var customerRecordBytes []byte
	var customerRecord *CustomerRecord
	var customerKey string
	var err error

	// Create the record key
	customerKey, err = GetCustomerRecordKey(stub, customerId)
	if err != nil {
		return customerRecord, err
	}

	// Pull the record bytes from the ledger
	customerRecordBytes, err = stub.GetState(customerKey)
	if err != nil {
		return customerRecord, err
	}

	if len(customerRecordBytes) == 0 {
		err = errors.New(fmt.Sprintf("No record found for Customer.ID %s", customerId))
		return customerRecord, err
	}

	// Unmarshal the JSON
	err = json.Unmarshal(customerRecordBytes, &customerRecord)
	if err != nil {
		return customerRecord, err
	}

	return customerRecord, nil
}

func SetCustomerRecord(stub shim.ChaincodeStubInterface, customerRecord *CustomerRecord) error {
	/*
		Sets a CustomerRecord object within the ledger

		Args:
			stub: HF shim interface
			customerRecord: CustomerRecord object to be set in the ledger

		Returns:
			err: Error object. nil if no error occurred.

	*/
	var customerRecordBytes []byte
	var customerKey string
	var err error

	// Create the record key
	customerKey, err = GetCustomerRecordKey(stub, customerRecord.Id)
	if err != nil {
		return err
	}

	// marshal the struct to JSON
	customerRecordBytes, err = json.Marshal(customerRecord)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshaling Customer record with CustomerRecord.ID %s",
			customerRecord.Id))
	}

	// Push the record back to the ledger
	err = stub.PutState(customerKey, customerRecordBytes)
	if err != nil {
		return err
	}

	return nil
}

func GetBankAccount(stub shim.ChaincodeStubInterface, bankAccountId string) (*BankAccount, error) {
	/*
		Fetches a BankAccount object from off the ledger

		Args:
			stub: HF shim interface
			bankAccountId: Primary Key of the BankAccount record

		Returns:
			bankAccount: BankAccount struct obj for the requested record
			err: Error object. nil if no error occurred.

	*/
	var bankAccountBytes []byte
	var bankAccount *BankAccount
	var bankAccountKey string
	var err error

	// Create the record key
	bankAccountKey, err = GetBankAccountKey(stub, bankAccountId)
	if err != nil {
		return bankAccount, err
	}

	// Pull the record bytes from the ledger
	bankAccountBytes, err = stub.GetState(bankAccountKey)
	if err != nil {
		return bankAccount, err
	}

	if len(bankAccountBytes) == 0 {
		err = errors.New(fmt.Sprintf("No record found for BankAccount.ID %s", bankAccountId))
		return bankAccount, err
	}

	// Unmarshal the JSON
	err = json.Unmarshal(bankAccountBytes, &bankAccount)
	if err != nil {
		return bankAccount, err
	}

	return bankAccount, nil

}

func SetBankAccount(stub shim.ChaincodeStubInterface, bankAccount *BankAccount) error {
	/*
		Sets a BankAccount object within the ledger

		Args:
			stub: HF shim interface
			bankAccount: BankAccount object to be set in the ledger

		Returns:
			err: Error object. nil if no error occurred.

	*/
	var bankAccountBytes []byte
	var bankAccountKey string
	var err error

	// Create the record key
	bankAccountKey, err = GetBankAccountKey(stub, bankAccount.Id)
	if err != nil {
		return err
	}

	// Validate balance
	if bankAccount.Balance < 0.0 {
		return errors.New(fmt.Sprintf("cannot update Bank Account balance of $%.2f; Balance must be >= $0.0",
			bankAccount.Balance))
	}

	// marshal the struct to JSON
	bankAccountBytes, err = json.Marshal(bankAccount)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshaling BankAccount record with BankAccount.ID %s",
			bankAccount.Id))
	}

	// Push the record back to the ledger
	err = stub.PutState(bankAccountKey, bankAccountBytes)
	if err != nil {
		return err
	}

	return nil
}

func GetAppDevRecord(stub shim.ChaincodeStubInterface, appDevId string) (*AppDevRecord, error) {
	/*
		Fetches a AppDevRecord object from off the ledger

		Args:
			stub: HF shim interface
			appDevId: Primary Key of the AppDevRecord object

		Returns:
			appDevRecord: AppDevRecord struct obj for the requested record
			err: Error object. nil if no error occurred.

	*/
	var appDevRecordBytes []byte
	var appDevRecord *AppDevRecord
	var appDevRecordKey string
	var err error

	// Create the record key
	appDevRecordKey, err = GetAppDevRecordKey(stub, appDevId)
	if err != nil {
		return appDevRecord, err
	}

	// Pull the record bytes from the ledger
	appDevRecordBytes, err = stub.GetState(appDevRecordKey)
	if err != nil {
		return appDevRecord, err
	}

	if len(appDevRecordBytes) == 0 {
		err = errors.New(fmt.Sprintf("No Bank Account record found for BankAccount.ID %s", appDevId))
		return appDevRecord, err
	}

	// Unmarshal the JSON
	err = json.Unmarshal(appDevRecordBytes, &appDevRecord)
	if err != nil {
		return appDevRecord, err
	}

	return appDevRecord, nil

}

func SetAppDevRecord(stub shim.ChaincodeStubInterface, appDevRecord *AppDevRecord) error {
	/*
		Sets an AppDevRecord object within the ledger

		Args:
			stub: HF shim interface
			appDevRecord: AppDevRecord object to be set in the ledger

		Returns:
			err: Error object. nil if no error occurred.

	*/
	var appDevRecordBytes []byte
	var appDevKey string
	var err error

	// Create the record key
	appDevKey, err = GetAppDevRecordKey(stub, appDevRecord.Id)
	if err != nil {
		return err
	}

	// marshal the struct to JSON
	appDevRecordBytes, err = json.Marshal(appDevRecord)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshaling AppDev record with AppDevRecord.ID %s",
			appDevRecord.Id))
	}

	// Push the record back to the ledger
	err = stub.PutState(appDevKey, appDevRecordBytes)
	if err != nil {
		return err
	}

	return nil
}
