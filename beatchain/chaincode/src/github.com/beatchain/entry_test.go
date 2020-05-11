package main

import (
	"github.com/beatchain/utils"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"testing"
)



func stringToBytes(strArray []string) [][]byte {
	var output [][]byte
	for _ , p := range strArray {
		output = append(output, []byte(p))
	}
	return output
}

func getInitArguments() [][]byte {
	return [][]byte{[]byte("init"),
		[]byte(utils.BEATCHAIN_ADMIN_BALANCE),
		[]byte(utils.TEST_APPDEV_ID),
		[]byte(utils.TEST_APPDEV_BA_ID),
		[]byte(utils.TEST_APPDEV_DEVSHARE),
		[]byte(utils.TEST_APPDEV_BA_BALANCE),
		[]byte(utils.TEST_CUSTOMER_ID),
		[]byte(utils.TEST_CUSTOMER_BA_ID),
		[]byte(utils.TEST_CUSTOMER_SUBFEE),
		[]byte(utils.TEST_CUSTOMER_SUB_DUE_DATE),
		[]byte(utils.TEST_CUSTOMER_BA_BALANCE),
		[]byte(utils.TEST_CREATOR_ID),
		[]byte(utils.TEST_CREATOR_BA_ID),
		[]byte(utils.TEST_CREATOR_BA_BALANCE),
		[]byte(utils.TEST_PRODUCT_ID),
		[]byte(utils.TEST_PRODUCT_NAME),
		[]byte(utils.TEST_PRODUCT_TOTAL_LISTENS),
		[]byte(utils.TEST_PRODUCT_UNREN_LISTENS),
		[]byte(utils.TEST_PRODUCT_TOTAL_METRICS),
		[]byte(utils.TEST_PRODUCT_UNREN_METRICS),
		[]byte(utils.TEST_PRODUCT_ADD_METRICS),
		[]byte(utils.TEST_PRODUCT_ACTIVE),
		[]byte(utils.TEST_CONTRACT_PPS),
		[]byte(utils.TEST_CONTRACT_STATUS)}
}

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
	bal, _ := strconv.ParseFloat(utils.BEATCHAIN_ADMIN_BALANCE, 32)
	utils.CheckBankAccount(t, stub, utils.BEATCHAIN_ADMIN_BANK_ACCOUNT_ID, float32(bal))
	bal, _ = strconv.ParseFloat(utils.TEST_APPDEV_BA_BALANCE, 32)
	utils.CheckBankAccount(t, stub, utils.TEST_APPDEV_BA_ID, float32(bal))
	bal, _ = strconv.ParseFloat(utils.TEST_CUSTOMER_BA_BALANCE, 32)
	utils.CheckBankAccount(t, stub, utils.TEST_CUSTOMER_BA_ID, float32(bal))
}

func beatchain_init(t *testing.T)  (*BeatchainChaincode, *shim.MockStub) {
	scc := new(BeatchainChaincode)
	scc.testMode = true
	stub := shim.NewMockStub("Beatchain", scc)
	checkInit(t, stub, getInitArguments())
	return scc, stub
}

func TestBeatchain_Init(t *testing.T) {
	beatchain_init(t)
}

func TestFullQueries(t *testing.T) {
	_, stub := beatchain_init(t)
	utils.ExecQuery(t, stub, "ListBankAccounts")
	utils.ExecQuery(t, stub, "ListAllCustomers")
}

func TestRenewal(t *testing.T) {
	_, stub := beatchain_init(t)
	res := stub.MockInvoke("1", [][]byte{[]byte("RenewSubscription")})
	if res.Status != shim.OK {
		fmt.Println("Query", "RenewSubscription", "failed", string(res.Message))
		t.FailNow()
	}
	utils.ExecQuery(t, stub, "ListBankAccounts")
	utils.ExecQuery(t, stub, "ListAllCustomers")
}

func TestCollection(t *testing.T) {
	_, stub := beatchain_init(t)
	utils.ExecQuery(t, stub, "CollectPayment")
	utils.ExecQuery(t, stub, "ListBankAccounts")
}

func TestTransferFunctions(t *testing.T) {
	_, stub := beatchain_init(t)
	amount64, err := strconv.ParseFloat(utils.TEST_APPDEV_BA_BALANCE, 64)
	if err != nil {
		fmt.Println("Failed to parse test appdev balance to float64:", utils.TEST_APPDEV_BA_BALANCE)
		t.FailNow()
	}
	// Test add money
	utils.ExecInvoke(t, stub, "TransferFunds", []string{utils.TEST_APPDEV_BA_ID, "100"})
	amount64 += 100.00
	utils.CheckBankAccount(t, stub, utils.TEST_APPDEV_BA_ID, float32(amount64))

	// Test withdrawal money
	utils.ExecInvoke(t, stub, "TransferFunds", []string{utils.TEST_APPDEV_BA_ID, "-100"})
	amount64 += -100.00
	utils.CheckBankAccount(t, stub, utils.TEST_APPDEV_BA_ID, float32(amount64))
}

func TestAddFunctions(t *testing.T) {
	var id *string
	_, stub := beatchain_init(t)

	// Add an appdev record
	id = utils.ExecInvoke(t, stub, "AddAppDevRecord", []string{"0.5"})
	fmt.Println("Returned id:", *id)
	appdevrec := utils.FetchTestAppdevRecord(t, stub, *id)
	fmt.Printf("Record: %+v\n", appdevrec)
	appdevbaRec := utils.FetchTestBankAccount(t, stub, appdevrec.BankAccountId)
	fmt.Printf("Bank Account Record: %+v\n", appdevbaRec)


	// Add a customer record
	id = utils.ExecInvoke(t, stub, "AddCustomerRecord", []string{"20.00"})
	fmt.Println("Returned id:", *id)
	custrec := utils.FetchTestCustomerRecord(t, stub, *id)
	fmt.Printf("Record: %+v\n", custrec)
	custbarec := utils.FetchTestBankAccount(t, stub, custrec.BankAccountId)
	fmt.Printf("Bank Account Record: %+v\n", custbarec)

	// Add a Creator record
	id = utils.ExecInvoke(t, stub, "AddCreatorRecord", []string{})
	fmt.Println("Returned id:", *id)
	createrec := utils.FetchTestCreatorRecord(t, stub, *id)
	fmt.Printf("Record: %+v\n", createrec)
	createbarec := utils.FetchTestBankAccount(t, stub, createrec.BankAccountId)
	fmt.Printf("Bank Account Record: %+v\n", createbarec)

	// Add a product record
	id = utils.ExecInvoke(t, stub, "AddProduct", []string{"test product name"})
	fmt.Println("Returned id:", *id)
	prodrec := utils.FetchTestProductRecord(t, stub, *id)
	fmt.Printf("Record: %+v\n", prodrec)

}


func TestContractFunctions(t *testing.T) {
	//var id *string
	var contract *utils.Contract
	_, stub := beatchain_init(t)

	_ = utils.ExecInvoke(t, stub, "OfferContract", []string{utils.TEST_APPDEV_ID, utils.TEST_CREATOR_ID, utils.TEST_PRODUCT_ID, "0.02"})
	contract = utils.FetchTestContractRecord(t, stub, utils.TEST_CREATOR_ID, utils.TEST_APPDEV_ID, utils.TEST_PRODUCT_ID)
	fmt.Printf("Contract: %+v\n", contract)

	_ = utils.ExecInvoke(t, stub, "AcceptContract", []string{utils.TEST_CREATOR_ID, utils.TEST_PRODUCT_ID, utils.TEST_APPDEV_ID})
	contract = utils.FetchTestContractRecord(t, stub, utils.TEST_CREATOR_ID, utils.TEST_APPDEV_ID, utils.TEST_PRODUCT_ID)
	fmt.Printf("Contract: %+v\n", contract)
}





