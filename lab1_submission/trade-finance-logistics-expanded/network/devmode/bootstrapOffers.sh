#!/usr/bin/env bash

peer chaincode install -p chaincodedev/chaincode/trade_workflow_v1 -n tw -v 0

# Setting up the accepted L/C prereq
peer chaincode instantiate -n tw -v 0 -c '{"Args":["init","LumberInc","LumberBank","100000", "WoodenToys","ToyBank","200000","UniversalFreight","ForestryDepartment"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["requestTrade", "foo", "70000", "Wood for Toys"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["acceptTrade", "foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getTradeStatus", "foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["requestLC", "foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getLCStatus", "foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["requestLC", "foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["issueLC", "foo", "lc8349", "12/31/2018", "E/L", "B/L"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["acceptLC", "foo"]}' -C tradechannel