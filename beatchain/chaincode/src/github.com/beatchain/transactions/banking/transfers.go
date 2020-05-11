package banking

import (
	"github.com/beatchain/utils"
	"fmt"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"math"
	"strconv"
)

func validateTransfer(transaction *utils.Transaction) (float32, error) {
	/*
		Validates the inputs to the TransferFunds function
	*/
	var amount float32
	var amount64 float64
	var err error

	// Access control: Only a Beatchain Admin Org member can invoke this transaction
	if !transaction.TestMode && !utils.AuthenticateBeatchainAdmin(transaction) {
		return amount, errors.New(fmt.Sprintf("caller not a member of Beatchain Admin Org. Access denied"))
	}
	// Validate no other args are specified
	if len(transaction.Args) != 2 {
		return amount, errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {BankAccountId, amount}. Found %d", len(transaction.Args)))
	}

	// Parse and validate amount
	amount64, err = strconv.ParseFloat(transaction.Args[1], 64)
	if err != nil {
		return amount, errors.New(fmt.Sprintf("Cannot parse amount to float64: %s", transaction.Args[1]))
	}
	amount64 = math.Round(amount64*100)/100
	if math.Abs(amount64) == 0.00 {
		return amount, errors.New(fmt.Sprintf("Cannot transfer amount of $0.00 (rounded)"))
	}
    // Limit the total amount in each txn
	if math.Abs(amount64) > 1000.0 {
		return amount, errors.New(fmt.Sprintf("Cannot transfer over $1000.00 in a single txn. Given: %.2f", amount64))
	}

	return float32(amount64), nil
}

func TransferFunds(stub shim.ChaincodeStubInterface, transaction *utils.Transaction) pb.Response {
	/*
		Credits monies transferred off-chain bank accounts to those on-chain (i.e. deposits and withdrawals)
	    Used by administrators to manually move funds in the ledger.

		Args:
			bankAccountId (string): ID of the BankAccount whose balance will be altered
	*/
	var bankAccountId string
	var bankAccount *utils.BankAccount
	var amount float32
	var err error

	// Validate inputs
	amount, err = validateTransfer(transaction)
	if err != nil {
		return shim.Error(err.Error())
	}
	bankAccountId = transaction.Args[0]

	// lookup Bank Account
	bankAccount, err = utils.GetBankAccount(stub, bankAccountId)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error accessing BA with id %s: %s", bankAccountId, err.Error()))
	}

	// Transfer and validate solvency
	bankAccount.Balance += amount
	if bankAccount.Balance < 0.00 {
		return shim.Error(fmt.Sprintf("BA ID: %s Insufficient Funds for payment of %.2f", bankAccountId, amount))
	}

	// Set change in ledger
	err = utils.SetBankAccount(stub, bankAccount)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("SUCCESS"))
}