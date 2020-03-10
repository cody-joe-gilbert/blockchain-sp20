#!/bin/bash
#
# Copyright 2018 IBM All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -e
cd /opt/trade/channel-artifacts

# This script expedites the chaincode development process by automating the
# requisite channel create/join commands

# We use a pre-generated orderer.block and channel transaction artifact (tradechannel.tx),
# both of which are created using the configtxgen tool

# first we create the channel against the specified configuration in tradechannel.tx
# this call returns a channel configuration block - myc.block - to the CLI container
peer channel create -c tradechannel -f ./channel.tx -o orderer:7050

# now we will join the channel and start the chain with tradechannel.block serving as the
# channel's first block (i.e. the genesis block)
peer channel join -b ./tradechannel.block

# Now the user can proceed to build and start chaincode in one terminal
# And leverage the CLI container to issue install instantiate invoke query commands in another


fabric-ca-client enroll -u http://admin:adminpw@ca:7054 --mspdir admin

ORG_NAME="importer"
fabric-ca-client register --id.name ${ORG_NAME} --id.secret pwd1 --id.type user \
    --id.attrs "tradelimit=1000:ecert,testorg=importer:ecert" -u http://ca:7054
fabric-ca-client enroll -u http://${ORG_NAME}:pwd1@ca:7054 \
    --enrollment.attrs "tradelimit,testorg,email:opt" --mspdir ${ORG_NAME}
mkdir ~/.fabric-ca-client/${ORG_NAME}/admincerts
cp -p ~/.fabric-ca-client/${ORG_NAME}/signcerts/*  ~/.fabric-ca-client/${ORG_NAME}/admincerts


ORG_NAME="exporter"
fabric-ca-client register --id.name ${ORG_NAME} --id.secret pwd1 --id.type user \
    --id.attrs "tradelimit=1000:ecert,testorg=exporter:ecert" -u http://ca:7054
fabric-ca-client enroll -u http://${ORG_NAME}:pwd1@ca:7054 \
    --enrollment.attrs "tradelimit,testorg,email:opt" --mspdir ${ORG_NAME}
mkdir ~/.fabric-ca-client/${ORG_NAME}/admincerts
cp -p ~/.fabric-ca-client/${ORG_NAME}/signcerts/*  ~/.fabric-ca-client/${ORG_NAME}/admincerts


ORG_NAME="lender"
fabric-ca-client register --id.name ${ORG_NAME} --id.secret pwd1 --id.type user \
    --id.attrs "tradelimit=1000:ecert,testorg=lender:ecert" -u http://ca:7054
fabric-ca-client enroll -u http://${ORG_NAME}:pwd1@ca:7054 \
    --enrollment.attrs "tradelimit,testorg,email:opt" --mspdir ${ORG_NAME}
mkdir ~/.fabric-ca-client/${ORG_NAME}/admincerts
cp -p ~/.fabric-ca-client/${ORG_NAME}/signcerts/*  ~/.fabric-ca-client/${ORG_NAME}/admincerts


peer chaincode install -p chaincodedev/chaincode/trade_workflow_v1 -n tw -v 0


