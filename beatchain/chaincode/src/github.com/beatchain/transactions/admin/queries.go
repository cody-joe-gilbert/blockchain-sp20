/*
Handles range queries on the ledger
Owner(s): Cody Gilbert
 */

package admin

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strings"

	"github.com/beatchain/utils"
)

func validateListBankAccounts(transaction *utils.Transaction) error {
	/*
		Validates the inputs to the ListBankAccounts function
	*/

	// Validate an ID is given
	if !transaction.TestMode && transaction.CreatorId == "" {
		return errors.New(fmt.Sprintf("calling user ID not found"))
	}
	// Validate no other args are specified
	if len(transaction.Args) != 0 {
		return errors.New(fmt.Sprintf("ListBankAccounts takes no arguments"))
	}
	return nil
}


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

	// Validate request
	err = validateListBankAccounts(transaction)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Create an iterator for fetching bank account keys
	keysIterator, err = stub.GetStateByPartialCompositeKey(utils.BANK_ACCOUNT_KEY_PREFIX, []string{})
	if err != nil {
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
		currentBankAccount, err = utils.GetBankAccount(stub, result.Key)
		if err != nil {
			// Errors print the current listing prior to the error for debug purposes
			jsonOutput = append(jsonOutput, fmt.Sprintf("keys operation failed. Error accessing Bank Account: %s", err))
			return shim.Error(strings.Join(jsonOutput, "\n"))
		}
		jsonOutput = append(jsonOutput, fmt.Sprintf("Bank Account ID: %s Balance: %.2f", currentBankAccount.Id, currentBankAccount.Balance))
	}
	resultMsg := strings.Join(jsonOutput, "\n")
	fmt.Print(resultMsg)

	return shim.Success([]byte(resultMsg))
}