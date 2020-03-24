# Beatchain Chaincode

### Structure

* `transactions`: **Owner: Everyone** Contains the Go code packages for processing blockchain transactions dispatched by `entry.go`.
* `vendor`: **Owner: Everyone** Third-party Go code packages
* `assets.go`: **Owner: Arun** Defines the asset data structures (e.g. Product, Account, Creator, etc.)
* `constants.go`: **Owner: Everyone** Defines constant-valued variables (i.e. keys)
* `entry.go`: **Owner: Julian** Defines the main Hyperledger Fabric `Init` and `Invoke` functions, and dispatches `Invoke` queries to 
    transactions defined in the `transactions` directory.
