package utils

import "github.com/hyperledger/fabric/core/chaincode/shim"

func GetCustomerRecordKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey(KEY_OBJECT_FORMAT, []string{CUSTOMER_RECORD_KEY_PREFIX, id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}

func GetCreatorRecordKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey(KEY_OBJECT_FORMAT, []string{CREATOR_RECORD_KEY_PREFIX, id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}

func GetContractKey(stub shim.ChaincodeStubInterface, creatorId string, appDevId string, productId string,) (string, error) {
	key, err := stub.CreateCompositeKey(KEY_OBJECT_FORMAT, []string{CONTRACT_KEY_PREFIX, creatorId, appDevId, productId})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}

func SplitContractKey(stub shim.ChaincodeStubInterface, key string) (string, string, string, error) {
	_, keyComponents, err := stub.SplitCompositeKey(key)
	if err != nil {
		return "", "", "", err
	}
	creatorId := keyComponents[1]
	appDevId := keyComponents[2]
	productId := keyComponents[3]
	return creatorId, appDevId, productId, nil

}

func GetBankAccountKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey(KEY_OBJECT_FORMAT, []string{BANK_ACCOUNT_KEY_PREFIX, id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}

func GetAppDevRecordKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey(KEY_OBJECT_FORMAT, []string{APPDEV_RECORD_KEY_PREFIX, id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}

func GetProductKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	key, err := stub.CreateCompositeKey(KEY_OBJECT_FORMAT, []string{PRODUCT_KEY_PREFIX, id})
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}