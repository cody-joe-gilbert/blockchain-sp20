# Transaction Chaincode

This directory contains the Go code used to define chaincode transactions. Transactions have been subdivided into 
separate packages according to their application areas.

### Structure

* `admin`: **Owner: Arun** defines the administrative transactions (e.g. Add a creator, Add a product, etc.)
* `banking`: **Owner: Cody** defines the banking and subscription management transactions 
    (e.g. Customer pays their subscription, Creator obtains payment, etc.)
* `streaming`: **Owner: Julian** defines the fundamental streaming and operation transactions 
    (e.g. Customer song requests, AppDev Stream validation, etc.)