# Beatchain: Music. Immutable.

This folder contains the application code for launching the Beatchain REST API server to process chaincode
transactions.

## Files
* `main.py`: Entry point for the `uvicorn` and `fastapi` webapp framework

## Folders
* `middleware`: Python package containing the function definitions used to interface the webapp in `main.py` to the
Hyperledger Fabric Beatchain network via the Fabric Python SDK APIs.

**Note:** The API current handles only a single query. These functions will be expanded upon as the chaincode develops.

### Installation
The Beatchain app requires Python Ver. 3.6+ to be installed on the system.

The Fabric Python SDK requires additional dev packages to be installed:
```shell script
sudo apt-get install python-dev python3-dev libssl-dev  # Ubuntu/Debian	
sudo yum install python-devel python3-devel openssl-devel  # Redhat/CentOS	
brew install python python3 openssl  # MacOS	
```

Once Python has been installed, create and activate a new virtual environment for Python as follows:
```shell script
python3 -m venv ~/beatchain_env 
source ~/beatchain_env/bin/activate
```

Install the required Python 3.6 dependencies as follows:
```shell script

pip3 install fastapi uvicorn
```

The Fabric Python SDK must be installed from its source files as follows:
```shell script
git clone https://github.com/hyperledger/fabric-sdk-py.git
cd fabric-sdk-py
make install
```


### Launching
Before launching the Beatchain application server, verify that the network has been setup. If using the Docker
network simulation, execute the following:
```shell script
cd ../network
./beatchain.sh down 
./beatchain.sh clean 
./beatchain.sh up
cd ../application
```
Validate that the Docker containers have been launched using the command `docker ps`. 
To launch the FastAPI server to receive API requests, launch the uvicorn ASGI server:

```shell script
uvicorn main:app --reload
```

The application will now handle API requests at `localhost:8000`. To see the API options available, view 
additional documentation, and run example queries a Swagger UI application can be accessed via browser at
```
http://127.0.0.1:8000/docs
```
or for a `ReDoc` version:
```
http://127.0.0.1:8000/redoc
```

Prior to executing chaincode `query` and `invoke` commands, the Beatchain channels and chaincode
must installed and initialized. To boostrap the network from the API, the `/admin/create_app` endpoint
may be used.