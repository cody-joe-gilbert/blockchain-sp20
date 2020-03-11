# Expanded Trade Workflow example
This directory holds the trade workflow example extension created by Cody Gilbert for Project 1 of the Blockchain and its Applications course.

## Problem
In the extended scenario, the Importer Bank does not have to pay the Exporter for 60 days after delivery of the shipment into the 
Importer’s country. 
Since the Exporter would prefer to have the money sooner, it will transfer ownership of the ”accepted L/C” to a 
third bank, the ”Lender”. When the Lender receives the accepted L/C it will transfer a discounted amount to the Exporter Bank. 
Later, it will request the full amount from the Importer Bank. To achieve this logic, you need to implement additional Organization type called 
Lender, create an additional asset type ”accepted L/C”, and create a few new transactions corresponding to creation/transfer of the accepted 
L/C, payment from the Lender to the Exporter Bank, and payment from the Importer Bank to the Lender.

### Design

#### Workflow
Once a L/C has been accepted by the Exporter organization, the Exporter may request a line of credit from a member of the
Lender organization, which if accepted will transfer the L/C beneficiary from the Exporter to the Lender in exchange for
the accepted line of credit. This takes place over the following steps:

1. Exporter requests a line of credit associated with the L/C
2. A Lender offers a discounted amount of the L/C in the line of credit
3. The Exporter approves of the discounted amount, transfers ownership of the L/C, and transfers the discounted amount 
from the Lender account to the Exporter Account
4. All other transactions take place normally; funds paid to the Exporter are instead paid to the Lender


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
5. `printCreditLine`: Queries the entire state of the CL
6. `printLC`: Queries the entire state of the letter of credit (validates ownership is transferred)

The following transactions were modified:
1. `requestPayment`: Can now be invoked by either the Exporter or the Lender if they own the L/C
2. `makePayment`: Now pays whomever holds the L/C beneficiary status

Requirements:
* ABAC-enforced functions
* The L/C must be in the 'ACCEPTED' status and good must be shipped prior to CL request
* The Lender-provided discounted amount must be greater than 0 and less than or equal to the remaining unpaid trade amount.
* The Lender cannot request payment if they do not ls
own the L/C

### Dev Example with CLI
This section covers how to execute the trade workflow extension using the development environment. The following
steps will execute the application with the extended orgs and transactions:

```shell script
# Starting in the networking folder

# Clean up the docker containers and files
./trade.sh down -d true
rm -rf ./devmode/channel-artifacts ./devmode/crypto-config ./devmode/logs ./devmode/docker-compose-e2e.yaml 

# Start the channel back up
./trade.sh up -d true

# Install the chaincode
docker exec -it chaincode bash
cd ./trade_workflow_v1; go build
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=tw:0 ./trade_workflow_v1

# In another terminal, start the cli for execution
docker exec -it cli bash

# Examples: Each of the following may be executed. Reset state after each execution.
# Example 1: Full end-to-end workflow
chmod +x /opt/trade/fullWorkflow.sh; /opt/trade/fullWorkflow.sh

# Example 2: Full end-to-end workflow without Lender Involvement
# Will raise an error when the lender attempts to request payment but doesn't hold the L/C
chmod +x /opt/trade/fullNoLending.sh; /opt/trade/fullNoLending.sh

# Example 3: Error Checks
chmod +x /opt/trade/setupChannel.sh; /opt/trade/setupChannel.sh
# Setup the execution until L/C is accepted
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

# Try to Request CL with the wrong credentials
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/importer
peer chaincode invoke -n tw -c '{"Args":["getCreditLine", "foo"]}' -C tradechannel  # FAIL
sleep 2
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/lender
peer chaincode invoke -n tw -c '{"Args":["getCreditLine", "foo"]}' -C tradechannel  # FAIL
sleep 2
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/exporter
peer chaincode invoke -n tw -c '{"Args":["getCreditLine", "foo"]}' -C tradechannel  # SUCCEED
sleep 2

# Offer with the wrong credentials
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/importer
peer chaincode invoke -n tw -c '{"Args":["offerCL", "foo", "4500"]}' -C tradechannel  # FAIL
sleep 2
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/exporter
peer chaincode invoke -n tw -c '{"Args":["offerCL", "foo", "4500"]}' -C tradechannel  # FAIL
sleep 2
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/lender
peer chaincode invoke -n tw -c '{"Args":["offerCL", "foo", "4500"]}' -C tradechannel  # SUCCEED
sleep 2

# Accept a CL with the wrong credentials
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/lender
peer chaincode invoke -n tw -c '{"Args":["acceptCL", "foo"]}' -C tradechannel  # FAIL
sleep 2
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/importer
peer chaincode invoke -n tw -c '{"Args":["acceptCL", "foo"]}' -C tradechannel  # FAIL
sleep 2
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/exporter
peer chaincode invoke -n tw -c '{"Args":["acceptCL", "foo"]}' -C tradechannel  # SUCCEED

```
### Example with Middleware

**NOTE: Due to errors in JS execution, this example is not functional at this time** 

This section covers how to execute the trade workflow extension using the middleware environment:

```shell script
# Starting in the networking folder

# Clean up the docker containers and files
./trade.sh down
./trade.sh stopneworg
./trade.sh clean

# Start the network with the new org
./trade.sh generate -c tradechannel
./trade.sh up

# Switch to the middleware folder and launch the original channel
cd ../middleware
node createTradeApp.js

# Switch back and launch the Lender containers
cd ../network
./trade.sh startneworg

# Switch to the middleware folder and upgrade the channel
cd ../middleware
node run-upgrade-channel.js
node new-org-join-channel.js 
node upgrade-chaincode.js

# Run the example execution
node runCreditScenario.js
```