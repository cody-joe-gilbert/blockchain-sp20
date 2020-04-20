package utils

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

func CheckBankAccount(t *testing.T, stub *shim.MockStub, id string, value float32) {
	var recordBytes []byte
	var record *BankAccount
	var key string
	var err error

	key, err = stub.CreateCompositeKey("object~id", []string{BANK_ACCOUNT_KEY_PREFIX, id})
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



func checkQuery(t *testing.T, stub *shim.MockStub, function string, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	payload := string(res.Payload)
	if payload != value {
		fmt.Println("Query value", name, "was", payload, "and not", value, "as expected")
		t.FailNow()
	}
}