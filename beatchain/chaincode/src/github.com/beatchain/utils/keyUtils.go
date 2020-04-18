package utils

import "github.com/hyperledger/fabric/core/chaincode/shim"

const CUSTOMER_RECORD_KEY_PREFIX = "CustomerRecord"
const CREATOR_RECORD_KEY_PREFIX = "CreatorRecord"
const CONTRACT_KEY_PREFIX = "Contract"
const BANK_ACCOUNT_KEY_PREFIX = "BankAccount"
const APPDEV_RECORD_KEY_PREFIX = "AppDevRecord"
const PRODUCT_KEY_PREFIX = "Product"

func GetCustomerRecordKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey(CUSTOMER_RECORD_KEY_PREFIX, []string{id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}

func GetCreatorRecordKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey(CREATOR_RECORD_KEY_PREFIX, []string{id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}

func GetContractKey(stub shim.ChaincodeStubInterface, creatorId string, appDevId string, productId string,) (string, error) {
	key, err := stub.CreateCompositeKey(CONTRACT_KEY_PREFIX, []string{creatorId, appDevId, productId})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}

func GetBankAccountKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey(BANK_ACCOUNT_KEY_PREFIX, []string{id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}

func GetAppDevRecordKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey(APPDEV_RECORD_KEY_PREFIX, []string{id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}

func GetProductKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey(PRODUCT_KEY_PREFIX, []string{id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}