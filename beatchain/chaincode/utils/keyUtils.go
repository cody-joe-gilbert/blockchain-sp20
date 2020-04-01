package utils

import "github.com/hyperledger/fabric/core/chaincode/shim"

func getTradeKey(stub shim.ChaincodeStubInterface, tradeID string) (string, error) {
	tradeKey, err := stub.CreateCompositeKey("Trade", []string{tradeID})
	if err != nil {
		return "", err
	} else {
		return tradeKey, nil
	}
}
