#!/usr/bin/env bash

# Install the middleware JS dependencies
cd ../middleware
npm install
cd ../network

# Setup the docker containers and files
./trade.sh down
./trade.sh clean
./trade.sh generate -c tradechannel
./trade.sh up

# Validate Docker containers
docker ps

# Invoke being the admin on the CLI
docker exec -ti cli bash


# Use middleware to setup the access priviledges
cd ../middleware
node createTradeApp.js


peer chaincode instantiate -n tw -v 0 -c '{"Args":["init","LumberInc","LumberBank","100000", "WoodenToys","ToyBank","200000","UniversalFreight","ForestryDepartment"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["requestTrade", "foo", "70000", "Wood for Toys"]}' -C tradechannel