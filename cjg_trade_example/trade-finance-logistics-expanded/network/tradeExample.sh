#!/usr/bin/env bash

# Install the middleware JS dependencies
cd ../middleware
npm install
cd ../network

# Setup the docker containers and files
./trade.sh down
./trade.sh clean
./trade.sh generate -c tradechannel


# Validate Docker containers
docker ps

docker exec -ti chaincode bash

# Use middleware to setup the access priviledges
cd ../middleware
node createTradeApp.js


# Running in Dev
./trade.sh down -d true
./trade.sh clean -d true
./trade.sh up -d true

docker exec -it chaincode bash
cd ./trade_workflow_v1
go build
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=tw:0 ./trade_workflow_v1

# Open new window
docker exec -it cli bash
set -e
cd /opt/trade/channel-artifacts
peer channel create -c tradechannel -f ./channel.tx -o orderer:7050
peer channel join -b ./tradechannel.block

chmod +x /opt/trade/createIdentity.sh
/opt/trade/createIdentity.sh




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

# Executing the new transactions
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/importer

# Request a CL
peer chaincode invoke -n tw -c '{"Args":["getCreditLine", "foo", "importer"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getCLStatus", "foo"]}' -C tradechannel

# Offer a CL
peer chaincode invoke -n tw -c '{"Args":["offerCL", "foo", "500"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getCLStatus", "foo"]}' -C tradechannel


# Accept a CL
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/importer
peer chaincode invoke -n tw -c '{"Args":["acceptCL", "foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getCLStatus", "foo"]}' -C tradechannel



