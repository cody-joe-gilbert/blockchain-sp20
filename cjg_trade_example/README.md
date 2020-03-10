# Expanded Trade Workflow example
This directory holds the trade workflow example extension created by Cody Gilbert for Project 1 of the Blockchain and its Applications course.

## Problem
In the extended scenario, the Importer Bank does not have to pay the Exporter for 60 days after delivery of the shipment into the 
Importer’s country. 

Since the Exporter would prefer to have the money sooner, it will transfer ownership of the ”accepted L/C” to a 
third bank, the ”Lender”. When the Lender receives the accepted L/C it will transfer a discounted amount to the Exporter Bank. 

Later, 
it will request the full amount from the Importer Bank. To achieve this logic, you need to implement additional Organization type called 
Lender, create an additional asset type ”accepted L/C”, and create a few new transactions corresponding to creation/transfer of the accepted 
L/C, payment from the Lender to the Exporter Bank, and payment from the Importer Bank to the Lender.

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
5. `printCreditLine`: Queries the entire state of the CL
6. `printLC`: Queries the entire state of the letter of credit (validates ownership is transferred)

Requirements:
* The L/C must be in the 'ACCEPTED' status prior to CL request
* The Exporter must accept the CL prior to the Importer, thus the Importer cannot accept the CL unless it is in the 'PENDING IMPORTER' status.
* The original L/C beneficiary must be changed to the Lender to allow the Importer to pay the Lender directly, and the Importer must
be notified of this change. This is why the Importer is required to approve the CL prior to 'ACCEPTED' status.
* The Lender-provided discounted amount must be greater than 0 and less than or equal to the original trade amount.

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

# Request a CL as the exporter
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/exporter  # Uses the exporter credentials
peer chaincode invoke -n tw -c '{"Args":["getCreditLine", "foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["printCreditLine", "foo"]}' -C tradechannel # Print the state

# The following will be rejected
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/lender
peer chaincode invoke -n tw -c '{"Args":["offerCL", "foo", "-500"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["offerCL", "foo", "50000000"]}' -C tradechannel

# Offer a CL of 4500 as the Lender
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/lender
peer chaincode invoke -n tw -c '{"Args":["offerCL", "foo", "4500"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["printCreditLine", "foo"]}' -C tradechannel # Print the state

# Accept a CL as the exporter
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/exporter
peer chaincode invoke -n tw -c '{"Args":["acceptCL", "foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["printCreditLine", "foo"]}' -C tradechannel # Print the state


# Importer makes the half-payment to the payee (exporter or lender)
peer chaincode invoke -n tw -c '{"Args":["requestPayment","foo"]}' -C tradechannel
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/importer
peer chaincode invoke -n tw -c '{"Args":["makePayment","foo"]}' -C tradechannel

# Progress shipment
peer chaincode invoke -n tw -c '{"Args":["acceptShipmentAndIssueBL","foo","fooBL","1/31/2030","JFK","EWR"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["updateShipmentLocation","foo","DESTINATION"]}' -C tradechannel

# Importer makes the rest of the payment
export CORE_PEER_MSPCONFIGPATH=/root/.fabric-ca-client/importer
peer chaincode invoke -n tw -c '{"Args":["requestPayment","foo"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["makePayment","foo"]}' -C tradechannel

# See the final state of the LC
peer chaincode invoke -n tw -c '{"Args":["printLC", "foo"]}' -C tradechannel

# Helpful
peer chaincode invoke -n tw -c '{"Args":["getAccountBalance","foo","exporter"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getAccountBalance","foo","importer"]}' -C tradechannel
peer chaincode invoke -n tw -c '{"Args":["getAccountBalance","foo","lender"]}' -C tradechannel
```

### Example with Middleware
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