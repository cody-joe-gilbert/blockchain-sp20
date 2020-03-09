# Expanded Trade Workflow example
This directory holds the trade workflow example extension created by Cody Gilbert for Project 1 of the Blockchain and its Applications course.

### Design

#### Workflow
Once a L/C has been accepted by the Exporter organization, the Exporter may request a line of credit from a member of the
Lender organization, which if accepted will transfer the L/C beneficiary from the Exporter to the Lender in exchange for
the accepted line of credit. This takes place over the following steps:

1. Exporter requests a line of credit associated with the L/C
2. A Lender offers a discounted amount of the L/C in the line of credit
3. The Exporter approves of the discounted amount and informs the Importer to transfer their payments (out of band)
4. The Importer formally approves the the credit line and changes the L/C beneficiary to the Lender

#### Implementation
To achieve the requirements of this project, the `Lender` organization was created to offer lines of credit to requesting 
`Exporter` organizations. To support these lines of credit, the `CreditLine` object was created 
(see `chaincode\src\github.com\trade_workflow_v1\assets.go`) to hold the names of the stakeholders, the offered credit line 
amount, and all other necessary attributes.
The following transactions were added to the workflow (see `chaincode\src\github.com\trade_workflow_v1\tradeWorkflow.go`):

1. `getCreditLine`: Executed by the Exporter with the name of the Importer to create a `CreditLine` object in the 'REQUESTED' status.
The name of the Importer is important to add to the request, as the transfer of L/C beneficiary will require the name of the Importer later. 
2. `offerCL`: Executed by the Lender with a given discounted amount to offer the line of credit to the Exporter. Changes the
`CreditLine` object to the 'OFFERED' status with a `DiscountAmount` attribute.
3. `acceptCL`: First executed by the Exporter to accept the terms of the credit line to change the `CreditLine` object to the 'PENDING IMPORTER' status.
Then executed by the Importer to confirm the transfer, shift the `CreditLine` object to the 'ACCEPTED' status and changes the original L/C object `Beneficiary` 
attribute to the Lender organization.
4. `getCLStatus`: Returns the status of the CL

Requirements:
* The L/C must be in the 'ACCEPTED' status prior to CL request
* The Exporter must accept the CL prior to the Importer, thus the Importer cannot accept the CL unless it is in the 'PENDING IMPORTER' status.
* The original L/C beneficiary must be changed to the Lender to allow the Importer to pay the Lender directly, and the Importer must
be notified of this change. This is why the Importer is required to approve the CL prior to 'ACCEPTED' status.

### Dev Example with CLI
This section covers how to execute the trade workflow extension using the development environment. The following
steps will execute the application with the extended orgs and transactions:

```shell script
# Starting in the networking folder

# Clean up the docker containers and files
./trade.sh down -d true
./trade.sh clean -d true
rm -rf ./devmode/channel-artifacts ./devmode/crypto-config ./devmode/logs ./devmode/docker-compose-e2e.yaml 

# Start the channel back up
./trade.sh up -d true

# Install the chaincode
docker exec -it chaincode bash
cd ./trade_workflow_v1
go build
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=tw:0 ./trade_workflow_v1

# In another terminal, start the cli for execution
docker exec -it cli bash

# Bootstrap the example: This script will create the channel, generate the org credentials,
# and run through the workflow until the L/C is in an ACCEPTED status.
chmod +x /opt/trade/setupChannel.sh
/opt/trade/setupChannel.sh

# See the current state of the LC
peer chaincode invoke -n tw -c '{"Args":["printLC", "foo"]}' -C tradechannel


peer chaincode invoke -n tw -c '{"Args":["printCreditLine", "foo"]}' -C tradechannel


# Request a CL as the exporter
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/exporter  # Uses the exporter credentials
peer chaincode invoke -n tw -c '{"Args":["getCreditLine", "foo", "importer"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getCLStatus", "foo"]}' -C tradechannel  # Checks the status

# Offer a CL of 500 as the Lender
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/lender
peer chaincode invoke -n tw -c '{"Args":["offerCL", "foo", "500"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getCLStatus", "foo"]}' -C tradechannel

# Accept a CL as the exporter
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/exporter
peer chaincode invoke -n tw -c '{"Args":["acceptCL", "foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getCLStatus", "foo"]}' -C tradechannel

# Accept a CL as the importer
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/importer
peer chaincode invoke -n tw -c '{"Args":["acceptCL", "foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getCLStatus", "foo"]}' -C tradechannel

```