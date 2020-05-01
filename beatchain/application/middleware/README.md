# Beatchain: Music. Immutable.

This Python package contains the function definitions used to interface the webapp in `main.py` to the
Hyperledger Fabric (HF) Beatchain network via the Fabric Python SDK APIs.

## Files
* `access_utils.py`: Defines functions for controlling user access to the HF network (i.e. user enrollment and registration)
* `constants.py`: Defines constant-valued parameters used in network configuration and chaincode management
* `create_app.py`: Defines functions used to bootstrap the network including channel creation, channel joining, chaincode
installation, and chaincode instantiation.
* `operations.py`: Defines functions for executing the chaincode (i.e. invocation and query transactions)