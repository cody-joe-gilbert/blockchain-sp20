# Beatchain: Music. Immutable.

This folder contains the chaincode used by the Beatchain application

### Structure

* `transactions`: Contains the Go code packages for processing blockchain transactions dispatched by `entry.go`.
* `utils`: Contains the Go code utility functions used by the main and transaction packages to factor out tedious operations.
* `vendor`: Third-party Go code packages
* `entry.go`: Defines the main Hyperledger Fabric `Init` and `Invoke` functions, and dispatches `Invoke` queries to 
    transactions defined in the `transactions` directory.
* `entry_test.go`: Contains chaincode testing framework using the HF testing environment.

