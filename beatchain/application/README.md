# Beatchain: Music. Immutable.

This folder contains the application code for launching the Beatchain REST API server to process chaincode
transactions.

**Note:** The API current handles only a single query. These functions will be expanded upon as the chaincode develops.

### Installation
Prior to executing the application, install the required Python 3.6 and Node.js dependencies as follows:

```shell script
pip3 install fastapi uvicorn Naked
cd ../middleware
npm install
cd ../application
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
Validate that the Docker containers have been launched using the command `docker ps`. Initialize the HF channels
and chaincode by launching the middleware code as follows:
```shell script
cd ../middleware
node createApp.js
cd ../application
```

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