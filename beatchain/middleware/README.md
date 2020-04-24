# Beatchain: Music. Immutable.

This folder contains the Hyperledger Fabric (HF) Javascript SDK middleware used to interface the application
REST API to the HF framework. The application itself is Python-based, and thus this folder contains several
Node.js wrappers executed directly by the Python code. These wrappers can be eliminated once a stable version
of the HF Python SDK is released.

### Configuration Files
The following files specify the configuration of the HF middleware, and must be updated with changes to 
the Beatchain network:

* `config.json`: Channel configuration file
* `constants.js`: Node.js file specifying channel and chaincode parameters

### HF Node.js SDK Scripts
The following scripts directly interface with the HF framework as middleware:

* `clientUtils`: User credential and CA management
* `create-channel`: Creates the channel as specified within `config.json` and `constants.js`
* `createApp.js`: Bootstraps the network creation by creating the Beatchain channels and installing and initializing
the chaincode.
* `install-chaincode.js`: Chaincode installation
* `instantiate-chaincode.js`: Chaincode instantiation
* `invoke-chaincode.js`: Chaincode invocation (doesn't return information)
* `join-channel.js`: Joins network peers to the channel
* `query-chaincode.js`: Chaincode querying (returns information)

### Python Adapter Scripts
The following scripts have been created to interface Python applications to the Node.js code:

* `query_wrapper.js`: Interfaces with `query-chaincode.js` to perform chaincode queries