#!/usr/bin/env bash

peer chaincode instantiate -n tw -v 0 -c '{"Args":["init","LumberInc","LumberBank","100000", "WoodenToys","ToyBank","200000","UniversalFreight","ForestryDepartment"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["requestTrade", "foo", "70000", "Wood for Toys"]}' -C tradechannel