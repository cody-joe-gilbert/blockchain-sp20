package banking

import (
	"blockchain-sp20/beatchain/chaincode/utils"
	"encoding/json"
	"fmt"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
    "time"
)

// Temporary definition until the final is added to assets.go
// TODO: replace with assets.go definition
type CustomerRecord struct {
	Id							string		`json:"id"`
	AppDevId					string		`json:"appdevid"`
	BankAccountId				string		`json:"bankaccountid"`
	SubscriptionDueDate			time.Time	`json:"subscriptionduedate"`
	QueuedSong					string		`json:"qeuedsong"`
	PreviousSong				string		`json:"previoussong"`
}

func renewSubscription(transaction utils.Transaction) pb.Response {
	/*
	Renews the Customer's subscription for a month. Transfers money from the Customer's
	bank account to the AppDev's account and extends the subscription due date by a month.

	Args:
		transaction: Creator's transaction info

	 */
	var tradeKey, lcKey string
	var tradeAgreementBytes, letterOfCreditBytes, exporterBytes []byte
	var customerRecord *CustomerRecord
	var err error

	// Access control: Only an Customer Org member can invoke this transaction
	if !utils.AuthenticateCustomer(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Customer Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Trade ID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// Lookup trade agreement from the ledger
	tradeKey, err = getTradeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	tradeAgreementBytes, err = stub.GetState(tradeKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(tradeAgreementBytes) == 0 {
		err = errors.New(fmt.Sprintf("No record found for trade ID %s", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify that the trade has been agreed to
	if tradeAgreement.Status != ACCEPTED {
		return shim.Error("Trade has not been accepted by the parties")
	}

	// Lookup exporter (L/C beneficiary)
	exporterBytes, err = stub.GetState(expKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	letterOfCredit = &LetterOfCredit{"", "", string(exporterBytes), tradeAgreement.Amount, []string{}, REQUESTED}
	letterOfCreditBytes, err = json.Marshal(letterOfCredit)
	if err != nil {
		return shim.Error("Error marshaling letter of credit structure")
	}

	// Write the state to the ledger
	lcKey, err = getLCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(lcKey, letterOfCreditBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Letter of Credit request for trade %s recorded\n", args[0])

	return shim.Success(nil)
}