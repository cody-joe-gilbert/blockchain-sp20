#!/bin/bash

chmod +x /opt/trade/setupChannel.sh
/opt/trade/setupChannel.sh

# Setting up the accepted L/C prereq
peer chaincode instantiate -n tw -v 0 -c '{"Args":["init","LumberInc","LumberBank","100000", "WoodenToys","ToyBank","200000","UniversalFreight","ForestryDepartment","LenderBros","300000"]}' -C tradechannel
sleep 3
peer chaincode invoke -n tw -c '{"Args":["requestTrade", "foo", "5000", "Wood for Toys"]}' -C tradechannel
sleep 2
peer chaincode invoke -n tw -c '{"Args":["acceptTrade", "foo"]}' -C tradechannel
sleep 2
peer chaincode invoke -n tw -c '{"Args":["requestLC", "foo"]}' -C tradechannel
sleep 2
peer chaincode invoke -n tw -c '{"Args":["issueLC", "foo", "fooLC", "12/31/2030", "E/L", "B/L"]}' -C tradechannel
sleep 2
peer chaincode invoke -n tw -c '{"Args":["acceptLC", "foo"]}' -C tradechannel
sleep 2
peer chaincode invoke -n tw -c '{"Args":["requestEL","foo"]}' -C tradechannel
sleep 2
peer chaincode invoke -n tw -c '{"Args":["issueEL","foo","fooLC","1/31/2030"]}' -C tradechannel
sleep 2
peer chaincode invoke -n tw -c '{"Args":["prepareShipment","foo"]}' -C tradechannel
sleep 2

#### Proceed without creating a CL

# Lender payment request will fail, as they don't own the L/C here
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/lender
peer chaincode invoke -n tw -c '{"Args":["requestPayment","foo"]}' -C tradechannel  # FAIL
sleep 2
# Importer makes the half-payment to the payee (exporter only)
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/exporter
peer chaincode invoke -n tw -c '{"Args":["requestPayment","foo"]}' -C tradechannel  # SUCCEED
sleep 2
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/importer
peer chaincode invoke -n tw -c '{"Args":["makePayment","foo"]}' -C tradechannel
sleep 2
# See that money transfers
peer chaincode invoke -n tw -c '{"Args":["getAccountBalance","foo","exporter"]}' -C tradechannel
sleep 2
peer chaincode invoke -n tw -c '{"Args":["getAccountBalance","foo","importer"]}' -C tradechannel
sleep 2
peer chaincode invoke -n tw -c '{"Args":["getAccountBalance","foo","lender"]}' -C tradechannel
sleep 2
# Progress shipment to importer
peer chaincode invoke -n tw -c '{"Args":["acceptShipmentAndIssueBL","foo","fooBL","1/31/2030","JFK","EWR"]}' -C tradechannel
sleep 2
peer chaincode invoke -n tw -c '{"Args":["updateShipmentLocation","foo","DESTINATION"]}' -C tradechannel
sleep 2
# Importer makes the rest of the payment
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/exporter
peer chaincode invoke -n tw -c '{"Args":["requestPayment","foo"]}' -C tradechannel
sleep 2
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/importer
peer chaincode invoke -n tw -c '{"Args":["makePayment","foo"]}' -C tradechannel
sleep 2
# See the final balances
peer chaincode invoke -n tw -c '{"Args":["getAccountBalance","foo","exporter"]}' -C tradechannel
sleep 2
peer chaincode invoke -n tw -c '{"Args":["getAccountBalance","foo","importer"]}' -C tradechannel
sleep 2
peer chaincode invoke -n tw -c '{"Args":["getAccountBalance","foo","lender"]}' -C tradechannel