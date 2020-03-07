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

chmod +x /opt/trade/createIdentity.sh
/opt/trade/createIdentity.sh

chmod +x /opt/trade/setupChannel.sh
/opt/trade/setupChannel.sh



# Executing the new transactions
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/importer

# Request a CL
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/exporter
peer chaincode invoke -n tw -c '{"Args":["getCreditLine", "foo", "importer"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getCLStatus", "foo"]}' -C tradechannel

# Offer a CL
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/lender
peer chaincode invoke -n tw -c '{"Args":["offerCL", "foo", "500"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getCLStatus", "foo"]}' -C tradechannel


# Accept a CL
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/exporter
peer chaincode invoke -n tw -c '{"Args":["acceptCL", "foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getCLStatus", "foo"]}' -C tradechannel

export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/importer
peer chaincode invoke -n tw -c '{"Args":["acceptCL", "foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getCLStatus", "foo"]}' -C tradechannel

peer chaincode invoke -n tw -c '{"Args":["getLCStatus", "foo"]}' -C tradechannel

