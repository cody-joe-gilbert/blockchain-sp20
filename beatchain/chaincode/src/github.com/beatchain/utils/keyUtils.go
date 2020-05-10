package utils

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)


func GetUniqueId(stub shim.ChaincodeStubInterface, txn *Transaction) (string, error) {
	var byteId []byte
	var strId string
	var intId int64
	var err error

	if txn.LastUniqueId < 0  {
		/*
		If this is the first time this function has been called in
		this instance of the chaincode, fetch the unique key from the ledger
		 */
		// Fetch the last unique ID
		fmt.Println("LastUnique < 0")
		byteId, err = stub.GetState(UNIQUE_ID_KEY)
		if err != nil {
			return "", err
		}

		// Convert to string
		strId = string(byteId)

		// convert to int
		intId, err = strconv.ParseInt(strId, 10, 64)
		if err != nil {
			return strId, err
		}
	} else {
		/*
		If GetUniqueId has been executed already in the chaincode,
		the unique ID fetched from the ledger will be stale.
		This section will use the last key in memory and
		update the ledger accordingly.
		 */
		fmt.Println("LastUnique > 0")
		intId = txn.LastUniqueId
		strId = strconv.FormatInt(intId, 10)

	}
	// increment to get a new unique ID
	intId += 1
	txn.LastUniqueId = intId

	// Set new unique ID back on the ledger
	err = stub.PutState(UNIQUE_ID_KEY, []byte(strconv.FormatInt(intId, 10)))
	if err != nil {
		return strId, err
	}

	return strId, nil
}


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