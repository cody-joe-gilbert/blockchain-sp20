package utils

import "github.com/hyperledger/fabric/core/chaincode/shim"

func GetCustomerRecordKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey("CustomerRecord", []string{id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}

func GetBankAccountKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey("BankAccount", []string{id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}

func GetAppDevRecordKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey("AppDevRecord", []string{id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}