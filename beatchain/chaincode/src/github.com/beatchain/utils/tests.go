package utils

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
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




func CheckProduct(t *testing.T, stub *shim.MockStub, id string, value float32) {
	var recordBytes []byte
	var record *Product
	var key string
	var err error

	key, err = stub.CreateCompositeKey("object~id", []string{PRODUCT_KEY_PREFIX, id})
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

}

func ExecQuery(t *testing.T, stub *shim.MockStub, function string) {
	fmt.Println("Executing Query function:", function)
	res := stub.MockInvoke("1", [][]byte{[]byte(function)})
	if res.Status != shim.OK {
		fmt.Println("Query", function, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", function, "failed to get value")
		t.FailNow()
	}
	payload := string(res.Payload)
	fmt.Println(payload)
}

func ExecInvoke(t *testing.T, stub *shim.MockStub, function string, args []string) *string {
	fmt.Println("Executing invoke function:", function)

	var byteArgs [][]byte
	byteArgs = append(byteArgs, []byte(function))
	for i, s := range args {
		fmt.Println("Arg:", i, "Value:", s)
		byteArgs = append(byteArgs, []byte(s))
	}

	res := stub.MockInvoke("1", byteArgs)
	if res.Status != shim.OK {
		fmt.Println("Invoke", function, "failed", string(res.Message))
		t.FailNow()
	}

	if res.Payload != nil {
		payload := string(res.Payload)
		fmt.Println("Returned payload:", payload)
		return &payload
	}
	return nil

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


func FetchTestBankAccount(t *testing.T, stub *shim.MockStub, id string,) *BankAccount {
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

	return record
}


func FetchTestAppdevRecord(t *testing.T, stub *shim.MockStub, id string,) *AppDevRecord {
	var recordBytes []byte
	var record *AppDevRecord
	var key string
	var err error

	key, err = stub.CreateCompositeKey("object~id", []string{APPDEV_RECORD_KEY_PREFIX, id})
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

	return record
}


func FetchTestCustomerRecord(t *testing.T, stub *shim.MockStub, id string,) *CustomerRecord {
	var recordBytes []byte
	var record *CustomerRecord
	var key string
	var err error

	key, err = stub.CreateCompositeKey("object~id", []string{CUSTOMER_RECORD_KEY_PREFIX, id})
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

	return record
}

func FetchTestCreatorRecord(t *testing.T, stub *shim.MockStub, id string,) *CreatorRecord {
	var recordBytes []byte
	var record *CreatorRecord
	var key string
	var err error

	key, err = stub.CreateCompositeKey("object~id", []string{CREATOR_RECORD_KEY_PREFIX, id})
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

	return record
}


func FetchTestProductRecord(t *testing.T, stub *shim.MockStub, id string,) *Product {
	var recordBytes []byte
	var record *Product
	var key string
	var err error

	key, err = stub.CreateCompositeKey("object~id", []string{PRODUCT_KEY_PREFIX, id})
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

	return record
}