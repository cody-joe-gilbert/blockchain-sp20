/*
Handles transactions to renew a Customer's subscription by a month
Owner(s): Cody Gilbert
*/
package banking

import (
	"github.com/beatchain/utils"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"math"
	"time"
)



func validateRenewSubscription(transaction *utils.Transaction) error {
	/*
	Validates the inputs to the renewSubscription function
	 */
	// Access control: Only an Customer Org member can invoke this transaction
	if !transaction.TestMode && !utils.AuthenticateCustomer(transaction) {
		return errors.New(fmt.Sprintf("caller not a member of Customer Org. Access denied"))
	}
	if transaction.TestMode {
		transaction.CreatorId = utils.TEST_CUSTOMER_ID
	}
	// Validate an ID is given
	if !transaction.TestMode && transaction.CreatorId == "" {
		return errors.New(fmt.Sprintf("customer ID not found"))
	}
	// Validate no other args are specified
	if len(transaction.Args) != 0 {
		return errors.New(fmt.Sprintf("renewSubscription takes no arguments"))
	}


	return nil
}

func RenewSubscription(stub shim.ChaincodeStubInterface, transaction *utils.Transaction) pb.Response {
	/*
	Renews the Customer's subscription for a month. Transfers money from the Customer's
	bank account to the AppDev's account and extends the subscription due date by a month.

	Args:
		transaction: Creator's transaction info

	 */
	var customerRecord *utils.CustomerRecord
	var customerBankAccount, appDevBankAccount, beatchainAdminBankAccount *utils.BankAccount
	var appDevRecord *utils.AppDevRecord
	var appDevShare float64
	var err error

	// Validate inputs
	err = validateRenewSubscription(transaction)
	if err != nil {
		return shim.Error(err.Error())
	}

	// lookup customer record
	customerRecord, err = utils.GetCustomerRecord(stub, transaction.CreatorId)
	if err != nil {
		return shim.Error(err.Error())
	}

	// lookup customer bank account balance
	customerBankAccount, err = utils.GetBankAccount(stub, customerRecord.BankAccountId)
	if err != nil {
		return shim.Error(err.Error())
	}

	// lookup AppDev record
	appDevRecord, err = utils.GetAppDevRecord(stub, customerRecord.AppDevId)
	if err != nil {
		return shim.Error(err.Error())
	}

	// lookup AppDev Bank Account
	appDevBankAccount, err = utils.GetBankAccount(stub, appDevRecord.BankAccountId)
	if err != nil {
		return shim.Error(err.Error())
	}

	// lookup Beatchain Admin Bank Account
	beatchainAdminBankAccount, err = utils.GetBankAccount(stub, utils.BEATCHAIN_ADMIN_BANK_ACCOUNT_ID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Validate the user can pay for the subscription
	if customerBankAccount.Balance < customerRecord.SubscriptionFee {
		err = errors.New(fmt.Sprintf("Bank Account Balance $%.2f insufficient for fee of $%.2f",
			customerBankAccount.Balance, customerRecord.SubscriptionFee))
		return shim.Error(err.Error())
	}

	// Exchange funds, taking care that cents are appropriately handled
	customerBankAccount.Balance -= customerRecord.SubscriptionFee
	appDevShare = float64(customerRecord.SubscriptionFee * (1. - appDevRecord.AdminFeeFrac))
	appDevShare = math.Round(appDevShare*100)/100
	appDevBankAccount.Balance += float32(appDevShare)
	beatchainAdminBankAccount.Balance += customerRecord.SubscriptionFee - float32(appDevShare)

	// Increment subscription time
	if customerRecord.SubscriptionDueDate.Before(time.Now()){
		// If subscription lapsed, add 30 days from now
		customerRecord.SubscriptionDueDate = time.Now().Add(time.Hour * 24 * 30)
	} else {
		// if due date hasn't passed, add 30 days to the existing due date
		customerRecord.SubscriptionDueDate = customerRecord.SubscriptionDueDate.Add(time.Hour * 24 * 30)
	}

	// Save the changes to the ledger
	err = utils.SetCustomerRecord(stub, customerRecord)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = utils.SetBankAccount(stub, customerBankAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = utils.SetBankAccount(stub, appDevBankAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = utils.SetBankAccount(stub, beatchainAdminBankAccount)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("SUCCESS"))
}