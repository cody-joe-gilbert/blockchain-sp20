package utils

import (
	"encoding/json"
	"fmt"
	"github.com/beatchain/utils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

func CheckBankAccount(t *testing.T, stub *shim.MockStub, id string, value float32) {
	var recordBytes []byte
	var record *utils.BankAccount
	var key string
	var err error

	key, err = stub.CreateCompositeKey(BANK_ACCOUNT_KEY_PREFIX, []string{id})
	if err != nil {
		fmt.Println("Cannot create key from id: ", id)
		t.FailNow()
	}

	recordBytes, err = stub.GetState(key)
	if err != nil {
		fmt.Println("Cannot find record for key: ", key)
		t.FailNow()
	}

	err = json.Unmarshal(recordBytes, &record)
	if err != nil {
		fmt.Println("Cannot unmarshal record for key: ", key)
		t.FailNow()
	}

	if record.Balance != value {
		fmt.Printf("BA Balance %.2f != %.2f expected for key %s", record.Balance, value, key)
		t.FailNow()
	}

}
