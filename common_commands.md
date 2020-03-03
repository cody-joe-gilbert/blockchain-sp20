# Hyperledger Fabric Tools

### Transaction Example

To start the transaction scripts:

```shell script
bash
cd /Users/codygilbert/Desktop/go/src/trade-finance-logistics/network
./trade.sh up -d true
```

To compile the chain code and run it in a container:
```shell script
docker exec -it chaincode bash
cd trade_workflow_v2
go build
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=tw:0 ./trade_workflow_v2
```


To install the chaincode on the peer and instantiate the trade
```shell script
docker exec -it cli bash
peer chaincode install -p chaincodedev/chaincode/trade_workflow_v2 -n tw -v 0
peer chaincode instantiate -n tw -v 0 -c '{"Args":["init","LumberInc","LumberBank","100000","WoodenToys","ToyBank","200000","UniversalFreight","ForestryDepartment"]}' -C tradechannel
```

Moving a package of chaincode to the peers and building it
```shell script
# Update the code
CODE_DIR=$1
FILE=$(basename $CODE_DIR)
docker cp $FILE chaincode:$CODE_DIR
```

To shut the network down:
```shell script
./trade.sh down –d true
```

To clean the existing ledgers:

```shell script
./trade.sh clean –d true
```



### Starting 
To start a bash shell for running HF:

```shell script
bash
```

Ensure that `.bashrc` contains the following:
```shell script
export GOPATH=$HOME/Desktop/go
export PATH=$PATH:$(go env GOPATH)/bin
export PATH="/usr/local/opt/gnutar/libexec/gnubin:$PATH"
export PATH="$PATH:$GOPATH/src/github.com/hyperledger/fabric/.build/bin"
```

# Middleware

To install the middleware dependencies, first move to the directory containing the node.js code
and execute
```shell script
npm install package.json
```

Execute the trade instance:
```shell script
node createTradeApp.js
```

