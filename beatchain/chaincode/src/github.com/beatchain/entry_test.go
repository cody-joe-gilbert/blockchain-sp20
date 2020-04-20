package main

import (
	"github.com/beatchain/utils"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"testing"
)


const BEATCHAIN_ADMIN_BALANCE = "1000"
const TEST_APPDEV_ID = "1111"
const TEST_APPDEV_BA_ID = "1111"
const TEST_APPDEV_DEVSHARE = "0.1"
const TEST_APPDEV_BA_BALANCE = "1000"
const TEST_CUSTOMER_ID = "2222"
const TEST_CUSTOMER_BA_ID = "1111"
const TEST_CUSTOMER_SUBFEE = "1.00"
const TEST_CUSTOMER_SUB_DUE_DATE = "2020-06-01"
const TEST_CUSTOMER_BA_BALANCE = "1000"


func getInitArguments() [][]byte {
	return [][]byte{[]byte("init"),
		[]byte(BEATCHAIN_ADMIN_BALANCE), // Beatchain admin BA balance
		[]byte(TEST_APPDEV_ID), // Test Appdev ID
		[]byte(TEST_APPDEV_BA_ID), // Test Appdev BA ID
		[]byte(TEST_APPDEV_DEVSHARE),  // Test AdminFeeFrac BA ID
		[]byte(TEST_APPDEV_BA_BALANCE), // Test AppDev BankAccount Initial Balance
		[]byte(TEST_CUSTOMER_ID), // Test Customer ID
		[]byte(TEST_CUSTOMER_BA_ID),  // Test Customer BankAccount ID
		[]byte(TEST_CUSTOMER_SUBFEE),  // Test Customer SubscriptionFee
		[]byte(TEST_CUSTOMER_SUB_DUE_DATE),  // Test Customer SubscriptionDueDate
		[]byte(TEST_CUSTOMER_BA_BALANCE)}  // Test Customer BankAccount Initial Balance
}

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}


func TestBeatchain_Init(t *testing.T) {
	scc := new(BeatchainChaincode)
	scc.testMode = true
	stub := shim.NewMockStub("Beatchain", scc)

	// Init
	scc.testCalledFunction = "init"
	scc.testCreatorId = "test"
	scc.testCreatorOrg = utils.APPDEV_MSP
	scc.testCreatorCertIssuer = utils.APPDEV_CA
	scc.testArgs = []string{"test"}
	checkInit(t, stub, getInitArguments())

	bal, _ := strconv.ParseFloat(BEATCHAIN_ADMIN_BALANCE, 32)
	utils.CheckBankAccount(t, stub, utils.BEATCHAIN_ADMIN_BANK_ACCOUNT_ID, float32(bal))
}
