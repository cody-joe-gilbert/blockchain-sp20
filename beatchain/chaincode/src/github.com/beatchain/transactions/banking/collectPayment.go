package banking

import (
	"errors"
	"fmt"
	"github.com/beatchain/utils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"math"
	"strings"
)

func validateCollectPayment(transaction *utils.Transaction) error {
	/*
		Validates the inputs to the renewSubscription function
	*/
	// Access control: Only an Customer Org member can invoke this transaction
	if !transaction.TestMode && !utils.AuthenticateCreator(transaction) {
		return errors.New(fmt.Sprintf("caller not a member of Creator Org. Access denied"))
	}
	if transaction.TestMode {
		transaction.CreatorId = utils.TEST_CREATOR_ID
	}
	// Validate an ID is given
	if transaction.CreatorId == "" {
		return errors.New(fmt.Sprintf("user ID not found"))
	}
	// Validate no other args are specified
	if len(transaction.Args) != 0 {
		return errors.New(fmt.Sprintf("CollectPayment takes no arguments"))
	}
	return nil
}

func CollectPayment(stub shim.ChaincodeStubInterface, transaction *utils.Transaction) pb.Response {
	/*
		Processes payment for a creator by accumulating all product streams and withdrawing payments from the
		AppDev accounts from whom the product was streamed.

		Args:
			transaction: Creator's transaction info

	*/
	var appDevRecord *utils.AppDevRecord
	var creatorRecord *utils.CreatorRecord
	var currentProduct *utils.Product
	var currentContract *utils.Contract
	var creatorBankAccount, appDevBankAccount *utils.BankAccount
	var keysIterator shim.StateQueryIteratorInterface
	var paymentExceptions int32
	var payment, totalPayment float32
	var currentAppDevId, currentProductId string
	var paymentDetails []string
	var err error

	// Validate inputs
	err = validateCollectPayment(transaction)
	if err != nil {
		return shim.Error(err.Error())
	}

	// lookup Creator's record
	creatorRecord, err = utils.GetCreatorRecord(stub, transaction.CreatorId)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error accessing creatorRecord with id %s: %s", transaction.CreatorId, err.Error()))
	}

	// lookup Creator's  Bank Account
	creatorBankAccount, err = utils.GetBankAccount(stub, transaction.CreatorId)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error accessing creatorRecord BA with id %s: %s", creatorRecord.BankAccountId, err.Error()))
	}

	totalPayment = 0.0
	paymentExceptions = 0

	// Create an iterator for fetching creator's contract keys
	keysIterator, err = stub.GetStateByPartialCompositeKey(utils.KEY_OBJECT_FORMAT, []string{utils.CONTRACT_KEY_PREFIX, transaction.CreatorId})
	if err != nil {
		fmt.Print("Key iterator error: ")
		return shim.Error(err.Error())
	}
	defer keysIterator.Close()

	// Loop through contracts and process payments
	for keysIterator.HasNext() {
		result, err := keysIterator.Next()
		if err != nil {
			// Errors print the current listing prior to the error for debug purposes
			return shim.Error(fmt.Sprintf("Contract iteration operation failed: %s", err.Error()))
		}

		// Split the key into appDevId and productId
		_, currentAppDevId, currentProductId, err = utils.SplitContractKey(stub, result.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		// Fetch the contract
		currentContract, err = utils.GetContract(stub, transaction.CreatorId, currentAppDevId, currentProductId)
		if err != nil {
			// Errors print the current listing prior to the error for debug purposes
			return shim.Error(fmt.Sprintf("Error accessing Contract with key %s: %s", result.Key, err.Error()))
		}
		// lookup AppDev record
		appDevRecord, err = utils.GetAppDevRecord(stub, currentAppDevId)
		if err != nil {
			return shim.Error(fmt.Sprintf("Error accessing appDevRecord with id %s: %s", currentAppDevId, err.Error()))
		}

		// lookup AppDev Bank Account
		appDevBankAccount, err = utils.GetBankAccount(stub, appDevRecord.BankAccountId)
		if err != nil {
			return shim.Error(fmt.Sprintf("Error accessing appDevRecord BA with id %s: %s", appDevRecord.BankAccountId, err.Error()))
		}

		// lookup product record
		currentProduct, err = utils.GetProduct(stub, currentProductId)
		if err != nil {
			return shim.Error(fmt.Sprintf("Error accessing product with id %s: %s", currentProductId, err.Error()))
		}
		if !currentProduct.IsActive {
			// Skip "deleted" products
			continue
		}

		// Attempt to transfer funds
		// Exchange funds, taking care that cents are appropriately handled
		payment64 := float64(currentProduct.UnRenumeratedListens) * float64(currentContract.CreatorPayPerStream)
		payment = float32(math.Round(payment64*100)/100)

		if payment == 0.0 {
			// No payment needed; skip processing
			continue
		}

		if appDevBankAccount.Balance < payment {
			// AppDev has insufficient funds to pay the creator; Note the exception to the user and continue
			paymentExceptions += 1
			msg := fmt.Sprintf(
				"WARNING! AppDev ID: %s Insufficient Funds for payment of %.2f in accordance with Contract %s",
				currentAppDevId, payment, result.Key)
			paymentDetails = append(paymentDetails, msg)
			continue
		}
		// If appDev has the funds, go ahead and process payment
		appDevBankAccount.Balance -= payment
		creatorBankAccount.Balance += payment
		totalPayment += payment

		// Print out the details for the payment
		msg := fmt.Sprintf(
			"Payment: $%.2f \n" +
				"\tAppDev ID: %s \n" +
				"\tNum. Streams: %d\n" +
				"\tPayment per Stream: $%.4f\n" +
				"\tIn accordance with Contract: %s",
			payment, currentAppDevId, currentProduct.UnRenumeratedListens, currentContract.CreatorPayPerStream, result.Key)
		paymentDetails = append(paymentDetails, msg)

		// Reset product metrics
		currentProduct.TotalListens += currentProduct.UnRenumeratedListens
		currentProduct.TotalMetrics += currentProduct.UnRenumeratedMetrics
		currentProduct.UnRenumeratedListens = 0
		currentProduct.UnRenumeratedMetrics = 0

		// Update changes ledger
		err = utils.SetProduct(stub, currentProduct)
		if err != nil {
			return shim.Error(err.Error())
		}
		err = utils.SetBankAccount(stub, appDevBankAccount)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	if totalPayment == 0 && paymentExceptions == 0 {
		// If there were no payments and no insufficient fund warnings, return with the message
		resultMsg := "No payable opportunities found."
		return shim.Success([]byte(resultMsg))
	} else if totalPayment == 0 && paymentExceptions != 0 {
		msg := fmt.Sprintf("No payments made. AppDevs found with insufficient funds")
		paymentDetails = append(paymentDetails, msg)
		resultMsg := strings.Join(paymentDetails, "\n")
		return shim.Success([]byte(resultMsg))
	} else {
		// Submit final payment to the creator on the ledger
		err = utils.SetBankAccount(stub, creatorBankAccount)
		if err != nil {
			return shim.Error(err.Error())
		}

		// Return final details message to the Creator
		paymentDetails = append(paymentDetails, fmt.Sprintf("Total Payment: %.2f", totalPayment))
		if paymentExceptions != 0 {
			paymentDetails = append(paymentDetails, fmt.Sprintf("WARNING: AppDevs found with insufficient funds"))
		}
		resultMsg := strings.Join(paymentDetails, "\n")
		return shim.Success([]byte(resultMsg))
	}


}
