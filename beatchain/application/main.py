# Blockchain & Applications Project 2: Beatchain
# Owner(s): Cody Gilbert

from fastapi import BackgroundTasks, FastAPI, Query
from fastapi.responses import JSONResponse
import middleware.constants as constants
import middleware.access_utils as access_utils
from middleware.create_app import create_app
import middleware.operations as operations


app = FastAPI()

### To simplify layout, all API entry points are included here despite the
# the length of this module.

##############################
## Admin Functions
##############################

@app.post('/admin/beatchain/create_app')
async def creation_request(req: constants.CreateAppRequest,
                           background_tasks: BackgroundTasks,
                           org_name: constants.OrgNames = Query(..., title="Organization Name"),
                           test_mode: bool = Query(False, title="Debug Initialization Mode Flag")):
    """
    Submits a request to bootstrap the application on the
    HF network.
    """
    try:
        background_tasks.add_task(create_app, org_name, req.admin_user_name, req.admin_password, test_mode)
    except Exception as e:
        content = {'Status': 'Application Creation Request failed',
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)
    content = {'Status': 'Application Creation Request Submitted. Note: App creation may still fail during initialization',
               'Error': None}
    return JSONResponse(status_code=201, content=content)

@app.get('/info/network_info')
async def network_info_request():
    """
    Returns the network configuration information
    """
    try:
        info = operations.get_network_info()
    except Exception as e:
        content = {'Status': 'Info Request failed',
                   'Response': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)
    content = {'Status': 'Info Request Succeeded',
               'Response': info,
               'Error': None}
    return JSONResponse(status_code=200, content=content)

@app.post('/info/channels')
async def channel_info_request(req: constants.CreateAppRequest,
                           org_name: constants.OrgNames = Query(..., title="Organization Name")):
    """
    Returns a listing of all active channels
    """
    try:
        info = await operations.get_channels(org_name, req.admin_user_name, req.admin_password)
    except Exception as e:
        content = {'Status': 'Channel Info Request Failed',
                   'Response': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)
    content = {'Status': 'Channel Info Request Succeeded',
               'Response': info,
               'Error': None}
    return JSONResponse(status_code=200, content=content)

@app.post('/info/block_info')
async def block_info_request(req: constants.CreateAppRequest,
                             org_name: constants.OrgNames = Query(..., title="Organization Name"),
                             channel_name: constants.ChannelNames = Query(constants.channel_name, title="Channel Name")):
    """
    Returns the information of the given channel's current block
    """
    try:
        info = await operations.get_block_info(org_name, req.admin_user_name, req.admin_password, channel_name)
    except Exception as e:
        content = {'Status': 'Block Info Request Failed',
                   'Response': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)
    content = {'Status': 'Block Info Request Succeeded',
               'Response': info,
               'Error': None}
    return JSONResponse(status_code=200, content=content)


@app.post('/info/instantiated_chaincodes')
async def inst_code_request(req: constants.CreateAppRequest,
                             org_name: constants.OrgNames = Query(..., title="Organization Name"),
                             channel_name: constants.ChannelNames = Query(constants.channel_name, title="Channel Name")):
    """
    Returns a listing of all the instantiated chaincodes within the given channel
    """
    try:
        info = await operations.get_instantiated_chaincodes(org_name, req.admin_user_name, req.admin_password, channel_name)
    except Exception as e:
        content = {'Status': 'Chaincode Info Request Failed',
                   'Response': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)
    content = {'Status': 'Chaincode Info Request Succeeded',
               'Response': info,
               'Error': None}
    return JSONResponse(status_code=200, content=content)


@app.post('/admin/beatchain/transfer')
async def transfer_request(req: constants.CreateAppRequest,
                           org_name: constants.OrgNames = Query(..., title="Organization Name"),
                           bank_account_id: str = Query(..., title="Bank Account ID"),
                           amount: float = Query(..., title="Transfer Amount (+/- for Deposite/Withdrawal)")
                           ):
    """
    Transfers funds in/out of a given Bank Account. The function serves as the main source
    and sink of monies on the ledger.

    Must be executed as a Beatchain Admin
    """
    try:
        amount = str(round(amount, 2))
        response = await operations.invoke(org_name,
                                           req.admin_user_name,
                                           req.admin_password,
                                           constants.channel_name,
                                           function='TransferFunds',
                                           args=[bank_account_id, amount])
    except Exception as e:
        content = {'Status': 'Transfer Request failed',
                   'Response': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)
    content = {'Status': 'Transfer Request Successful',
               'Response': response,
               'Error': None}
    return JSONResponse(status_code=200, content=content)

@app.post('/admin/beatchain/register')
async def register(req: constants.RegisterUserRequest,
             org_name: constants.OrgNames = Query(..., title="Organization Name"),
             ):
    """
    Submits a request to register a user on HF network using the given
    org's certificate authority. If successful, return's the user's
    login password with the response.
    Note: Passing a user secret via API is *not* secure, and this endpoint should be
    used for demo purposes only.
    """
    # TODO: Passing a secret back is NOT secure! This section is for demo only!
    try:
        secret = await access_utils.register_user(org_name, req)
    except Exception as e:
        content = {'Status': 'Registration Request Failed',
                   'Secret': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)
    content = {'Status': 'Registration Request Succeeded',
               'Secret': secret,
               'Error': None}
    return JSONResponse(status_code=201, content=content)


@app.post('/admin/creator/add_creator')
async def add_creator(req: constants.AddUserRecordRequest):
    """
    FUNCTION NOT YET IMPLEMENTED

    Adds a creator user to the CreatorOrg (creatororg.beatchain.com) and generates
    a CreatorRecord asset on the ledger.

    Returns the new creator's ledger ID and login secret.

    Request must be submitted by an admin of the CreatorOrg (creatororg.beatchain.com)
    """
    # TODO: Passing a secret back is NOT secure! This section is for demo only!
    response = None
    try:
        # First add creator to the ledger
        response = await operations.invoke('creatororg.beatchain.com',
                                           req.user_name,
                                           req.user_password,
                                           constants.channel_name,
                                           function='AddCreatorRecord',
                                           args=[])
    except Exception as e:
        content = {'Status': 'Failed to add creator to ledger',
                   'Creator ID': None,
                   'Creator Secret': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)

    # Attempt to parse out the creator ID from the response
    try:
        creator_id = int(response.split()[-1])
    except Exception as e:
        content = {'Status': 'Cannot parse int creator id from response: ' + response,
                   'Creator ID': None,
                   'Creator Secret': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)

    # Register the new creator user
    try:
        register_req = constants.RegisterUserRequest(
            admin_user_name=req.admin_user_name,
            admin_password=req.admin_password,
            user_name=req.user_name,
            user_password=req.user_password,
            role='client',
            attrs=[{'creator_id': str(creator_id)}])
        secret = await access_utils.register_user('creatororg.beatchain.com',
                                                  register_req)
    except Exception as e:
        content = {'Status': 'Creator User Creation Failed',
                   'Creator ID': creator_id,
                   'Creator Secret': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)

    content = {'Status': 'Creator Creation Request Successful',
               'Creator ID': creator_id,
               'Creator Secret': secret,
               'Error': None}
    return JSONResponse(status_code=201, content=content)


@app.post('/admin/appdev/add_appdev')
async def add_appdev(req: constants.AddUserRecordRequest,
                     admin_fee_frac: float = Query(..., title="Subscription Administration Fee Fraction")
                     ):
    """
    FUNCTION NOT YET IMPLEMENTED

    Adds an appdev user to the AppdevOrg (appdevorg.beatchain.com) and generates
    a AppDevRecord asset on the ledger.

    Returns the new appdev's ledger ID and login secret.

    Request must be submitted by an admin of the AppDevOrg (appdevorg.beatchain.com)
    """
    # TODO: Passing a secret back is NOT secure! This section is for demo only!
    response = None
    try:
        # First add appdev to the ledger
        response = await operations.invoke('appdevorg.beatchain.com',
                                           req.user_name,
                                           req.user_password,
                                           constants.channel_name,
                                           function='AddAppDevRecord',
                                           args=[str(round(admin_fee_frac, 3))])
    except Exception as e:
        content = {'Status': 'Failed to add appdev record to ledger',
                   'AppDev ID': None,
                   'AppDev Secret': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)

    # Attempt to parse out the appdev ID from the response
    try:
        appdev_id = int(response.split()[-1])
    except Exception as e:
        content = {'Status': 'Cannot parse int appdev_id from response: ' + response,
                   'Creator ID': None,
                   'Creator Secret': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)

    # Register the new creator user
    try:
        register_req = constants.RegisterUserRequest(
            admin_user_name=req.admin_user_name,
            admin_password=req.admin_password,
            user_name=req.user_name,
            user_password=req.user_password,
            role='client',
            attrs=[{'appdev_id': str(appdev_id)}])
        secret = await access_utils.register_user('appdevorg.beatchain.com',
                                                  register_req)
    except Exception as e:
        content = {'Status': 'AppDev User Creation Failed',
                   'AppDev ID': appdev_id,
                   'AppDev Secret': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)

    content = {'Status': 'AppDev Creation Request Successful',
               'AppDev ID': appdev_id,
               'AppDev Secret': secret,
               'Error': None}
    return JSONResponse(status_code=201, content=content)


@app.post('/admin/customer/add_customer')
async def add_customer(req: constants.AddUserRecordRequest,
                      appdev_id: int = Query(..., title="Customer's Subscribing AppDev ID")
                      ):
    """
    FUNCTION NOT YET IMPLEMENTED

    Adds a customer user to the CustomerOrg (customerorg.beatchain.com) and generates
    a CustomerRecord asset on the ledger. Note that the customer's subscription fee and
    due date are set separately by an AppDevOrg client.

    Returns the new customer's ledger ID and login secret.

    Request must be submitted by an admin of the CustomerOrg (customerorg.beatchain.com)
    """
    # TODO: Passing a secret back is NOT secure! This section is for demo only!
    response = None
    try:
        # First add customer to the ledger
        response = await operations.invoke('customerorg.beatchain.com',
                                           req.user_name,
                                           req.user_password,
                                           constants.channel_name,
                                           function='AddCustomerRecord',
                                           args=[str(appdev_id)])
    except Exception as e:
        content = {'Status': 'Failed to add Customer to ledger',
                   'Customer ID': None,
                   'Customer Secret': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)

    # Attempt to parse out the creator ID from the response
    try:
        customer_id = int(response.split()[-1])
    except Exception as e:
        content = {'Status': 'Cannot parse int Customer id from response: ' + response,
                   'Customer ID': None,
                   'Customer Secret': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)

    # Register the new creator user
    try:
        register_req = constants.RegisterUserRequest(
            admin_user_name=req.admin_user_name,
            admin_password=req.admin_password,
            user_name=req.user_name,
            user_password=req.user_password,
            role='client',
            attrs=[{'customer_id': str(customer_id)}])
        secret = await access_utils.register_user('customerorg.beatchain.com',
                                                  register_req)
    except Exception as e:
        content = {'Status': 'Customer User Creation Failed',
                   'Customer ID': customer_id,
                   'Customer Secret': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)

    content = {'Status': 'Customer Creation Request Successful',
               'Customer ID': customer_id,
               'Customer Secret': secret,
               'Error': None}
    return JSONResponse(status_code=201, content=content)


@app.post('/admin/add_product')
async def add_product(req: constants.AddProductRequest):
    """
    FUNCTION NOT YET IMPLEMENTED

    Submits a request to add a product to the ledger with the given
    product_name under the CreatorID of the creator that submits the
    product creation request. The results will return the chaincode-generated Product
    ID.

    Request must be submitted by the CreatorOrg (creatororg.beatchain.com)
    """

    try:
        response = await operations.invoke('creatororg.beatchain.com',
                                           req.user_name,
                                           req.user_password,
                                           constants.channel_name,
                                           function='AddProduct',
                                           args=[req.product_name])
    except Exception as e:
        content = {'Status': 'Product Creation failed',
                   'Product ID': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)
    content = {'Status': 'Product Creation Succeeded',
               'Product ID': response,
               'Error': repr(e)}
    return JSONResponse(status_code=201, content=content)


##############################
## Invoke Functions
##############################

@app.post('/invoke')
async def invoke_request(req: constants.InvokeRequest,
                   org_name: constants.OrgNames = Query(..., title="Organization Name"),
                   channel_name: constants.ChannelNames = Query(..., title="Network Channel Name"),
                   function: constants.InvokeFunctions = Query(..., title="Chaincode Function"),
                   ):
    """
    Submits a blockchain transaction invocation to a subset of the peers in
    the network.
    """
    try:
        response = await operations.invoke(org_name,
                                     req.user_name,
                                     req.user_password,
                                     channel_name,
                                     function,
                                     req.args)
    except Exception as e:
        content = {'Status': 'Invoke Request failed',
                   'Response': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)
    content = {'Status': 'Invoke Request successful',
               'Response': response,
               'Error': None}
    return JSONResponse(status_code=200, content=content)

@app.post('/query')
async def query_request(req: constants.InvokeRequest,
                  org_name: constants.OrgNames = Query(..., title="Organization Name"),
                  channel_name: constants.ChannelNames = Query(..., title="Network Channel Name"),
                  function: constants.QueryFunctions = Query(..., title="Chaincode Function"),
                  ):
    """
    Submits a ledger query to a single peer within the specified org.
    Note that queries will NOT submit any changes to the ledger state, and will ONLY
    reflect the ledger state of the peer to which the request was submitted.
    Therefore, this query does NOT guarantee consistency like invoke operations.
    """
    try:
        response = await operations.query(org_name,
                                     req.user_name,
                                     req.user_password,
                                     channel_name,
                                     function,
                                     req.args)
    except Exception as e:
        content = {'Status': 'Query Request Failed',
                   'Response': None,
                   'Error': repr(e)}
        return JSONResponse(status_code=500, content=content)
    content = {'Status': 'Query Request Successful',
               'Response': response,
               'Error': None}
    return JSONResponse(status_code=200, content=content)
