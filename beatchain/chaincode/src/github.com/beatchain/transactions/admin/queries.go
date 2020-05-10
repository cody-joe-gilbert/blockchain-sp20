/*
Handles range queries on the ledger
Owner(s): Cody Gilbert
 */

package admin

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strings"

	"github.com/beatchain/utils"
)


func ListBankAccounts(stub shim.ChaincodeStubInterface, transaction *utils.Transaction) pb.Response {
	/*
		Lists all of the bank accounts and their balances from the ledger

		Args:
			transaction: Creator's transaction info

	*/
	var currentBankAccount *utils.BankAccount

	var jsonOutput []string
	var err error
	var keysIterator shim.StateQueryIteratorInterface


	// Validate an ID is given
	if !transaction.TestMode && transaction.CreatorId == "" {
		return shim.Error(fmt.Sprintf("calling user ID not found"))
	}
	// Validate no other args are specified
	if len(transaction.Args) != 0 {
		return shim.Error(fmt.Sprintf("ListBankAccounts takes no arguments"))
	}

	// Create an iterator for fetching bank account keys
	keysIterator, err = stub.GetStateByPartialCompositeKey("object~id", []string{utils.BANK_ACCOUNT_KEY_PREFIX})
	if err != nil {
		fmt.Print("Key iterator error: ")
		return shim.Error(err.Error())
	}
	defer keysIterator.Close()

	// Loop through keys and print account balances
	for keysIterator.HasNext() {
		result, err := keysIterator.Next()
		if err != nil {
			// Errors print the current listing prior to the error for debug purposes
			jsonOutput = append(jsonOutput, fmt.Sprintf("keys operation failed. Error accessing state: %s", err))
			return shim.Error(strings.Join(jsonOutput, "\n"))
		}
		_, keyComponents, err := stub.SplitCompositeKey(result.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		currentBankAccount, err = utils.GetBankAccount(stub, keyComponents[1])
		if err != nil {
			// Errors print the current listing prior to the error for debug purposes
			jsonOutput = append(jsonOutput, fmt.Sprintf("keys operation failed. Error accessing Bank Account: %s", err))
			return shim.Error(strings.Join(jsonOutput, "\n"))
		}
		jsonOutput = append(jsonOutput, fmt.Sprintf("Bank Account ID: %s Balance: %.2f", currentBankAccount.Id, currentBankAccount.Balance))
	}
	resultMsg := strings.Join(jsonOutput, "\n")
	return shim.Success([]byte(resultMsg))
}


func ListAllCustomers(stub shim.ChaincodeStubInterface, transaction *utils.Transaction) pb.Response {
	/*
		Lists all of the customer records and their details from the ledger

		Args:
			transaction: Creator's transaction info

	*/
	var currentCustomerRecord *utils.CustomerRecord

	var jsonOutput []string
	var err error
	var keysIterator shim.StateQueryIteratorInterface


	// Validate an ID is given
	if !transaction.TestMode && transaction.CreatorId == "" {
		return shim.Error(fmt.Sprintf("calling user ID not found"))
	}
	// Validate no other args are specified
	if len(transaction.Args) != 0 {
		return shim.Error(fmt.Sprintf("ListCustomers takes no arguments"))
	}

	// Create an iterator for fetching keys
	keysIterator, err = stub.GetStateByPartialCompositeKey("object~id", []string{utils.CUSTOMER_RECORD_KEY_PREFIX})
	if err != nil {
		fmt.Print("Key iterator error: ")
		return shim.Error(err.Error())
	}
	defer keysIterator.Close()

	// Loop through keys and print account balances
	for keysIterator.HasNext() {
		result, err := keysIterator.Next()
		if err != nil {
			// Errors print the current listing prior to the error for debug purposes
			jsonOutput = append(jsonOutput, fmt.Sprintf("keys operation failed. Error accessing state: %s", err))
			return shim.Error(strings.Join(jsonOutput, "\n"))
		}
		_, keyComponents, err := stub.SplitCompositeKey(result.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		currentCustomerRecord, err = utils.GetCustomerRecord(stub, keyComponents[1])
		if err != nil {
			// Errors print the current listing prior to the error for debug purposes
			jsonOutput = append(jsonOutput, fmt.Sprintf("keys operation failed. Error accessing Bank Account: %s", err))
			return shim.Error(strings.Join(jsonOutput, "\n"))
		}
		msg := fmt.Sprintf(
			"Customer ID: %s \n" +
				"\tAppDevId: %s\n" +
				"\tBankAccountId: %s\n" +
				"\tSubscriptionFee: %0.2f\n" +
				"\tSubscriptionDueDate: %s\n" +
				"\tQueuedSong: %s\n" +
				"\tPreviousSong: %s",
			currentCustomerRecord.Id,
			currentCustomerRecord.AppDevId,
			currentCustomerRecord.BankAccountId,
			currentCustomerRecord.SubscriptionFee,
			currentCustomerRecord.SubscriptionDueDate.String(),
			currentCustomerRecord.QueuedSong,
			currentCustomerRecord.PreviousSong)
		jsonOutput = append(jsonOutput, msg)
	}
	resultMsg := strings.Join(jsonOutput, "\n")
	return shim.Success([]byte(resultMsg))
}

func ListAppCustomers(stub shim.ChaincodeStubInterface, transaction *utils.Transaction) pb.Response {
	/*
		Lists all of the customer records and their details from the ledger for a particular app dev eg: spotify

		Args:
			transaction: Creator's transaction info, AppdevID

	*/
	var currentCustomerRecord *utils.CustomerRecord

	var jsonOutput []string
	var err error
	var keysIterator shim.StateQueryIteratorInterface


	// Validate an ID is given
	if !transaction.TestMode && transaction.CreatorId == "" {
		return shim.Error(fmt.Sprintf("calling user ID not found"))
	}
	// Validate no other args are specified
	if len(transaction.Args) != 1 {
		return shim.Error(fmt.Sprintf("ListCustomers takes 1 argument : {AppdevId}"))
	}

	appDevId := transaction.Args[0]
	// check for valid AppDev
	_, err = utils.GetAppDevRecord(stub, appDevId)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Create an iterator for fetching keys
	keysIterator, err = stub.GetStateByPartialCompositeKey("object~id", []string{utils.CUSTOMER_RECORD_KEY_PREFIX})
	if err != nil {
		fmt.Print("Key iterator error: ")
		return shim.Error(err.Error())
	}
	defer keysIterator.Close()

	// Loop through keys and print account balances
	for keysIterator.HasNext() {
		result, err := keysIterator.Next()
		if err != nil {
			// Errors print the current listing prior to the error for debug purposes
			jsonOutput = append(jsonOutput, fmt.Sprintf("keys operation failed. Error accessing state: %s", err))
			return shim.Error(strings.Join(jsonOutput, "\n"))
		}
		_, keyComponents, err := stub.SplitCompositeKey(result.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		currentCustomerRecord, err = utils.GetCustomerRecord(stub, keyComponents[1])
		if err != nil {
			// Errors print the current listing prior to the error for debug purposes
			jsonOutput = append(jsonOutput, fmt.Sprintf("keys operation failed. Error accessing Bank Account: %s", err))
			return shim.Error(strings.Join(jsonOutput, "\n"))
		}
		msg := fmt.Sprintf(
			"Customer ID: %s \n" +
				"\tAppDevId: %s\n" +
				"\tBankAccountId: %s\n" +
				"\tSubscriptionFee: %0.2f\n" +
				"\tSubscriptionDueDate: %s\n" +
				"\tQueuedSong: %s\n" +
				"\tPreviousSong: %s",
			currentCustomerRecord.Id,
			currentCustomerRecord.AppDevId,
			currentCustomerRecord.BankAccountId,
			currentCustomerRecord.SubscriptionFee,
			currentCustomerRecord.SubscriptionDueDate.String(),
			currentCustomerRecord.QueuedSong,
			currentCustomerRecord.PreviousSong)
		
		if appDevId == currentCustomerRecord.AppDevId {
			jsonOutput = append(jsonOutput, msg)
		}

		
	}
	resultMsg := strings.Join(jsonOutput, "\n")
	return shim.Success([]byte(resultMsg))
}