# Blockchain & Applications Project 2: Beatchain
# Owner(s): Cody Gilbert

from fastapi import FastAPI, Query
import middleware.constants as constants
import middleware.access_utils as access_utils
from middleware.create_app import create_app
import middleware.operations as operations

app = FastAPI()

@app.post('/admin/create_app')
async def creation_request(req: constants.CreateAppRequest,
                     org_name: constants.OrgNames = Query(..., title="Organization Name"),
                     test_mode: bool = Query(False, title="Debug Initialization Mode Flag"),
                     ):
    """
    Submits a request to bootstrap the application on the
    HF network.
    """
    try:
        await create_app(org_name,
                   req.admin_user_name,
                   req.admin_password,
                   test_mode)
    except Exception as e:
        return {'Status': 'Application Creation Request failed',
                'Error': repr(e)}
    return {'Status': 'Application Created',
            'Error': None}

@app.post('/admin/register')
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
        return {'Status': 'Registration Request failed',
                'Secret': None,
                'Error': repr(e)}
    return {'Status': 'Registration Request Succeeded',
            'Secret': secret,
            'Error': None}

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
        return {'Status': 'Invoke Request failed',
                'Response': None,
                'Error': repr(e),
                }
    return {'Status': 'Invoke Request successful',
            'Response': response,
            'Error': None,
            }

@app.post('/query')
async def query_request(req: constants.InvokeRequest,
                  org_name: constants.OrgNames = Query(..., title="Organization Name"),
                  channel_name: constants.ChannelNames = Query(..., title="Network Channel Name"),
                  function: constants.QueryFunctions = Query(..., title="Chaincode Function"),
                  ):
    """
    Submits a ledger query to a single peer within the specified org.
    Note that queries will NOT submit any changes to the ledger state
    """
    try:
        response = await operations.query(org_name,
                                     req.user_name,
                                     req.user_password,
                                     channel_name,
                                     function,
                                     req.args)
    except Exception as e:
        return {'Status': 'Query Request failed',
                'Response': None,
                'Error': repr(e),
                }
    return {'Status': 'Invoke Request successful',
            'Response': response,
            'Error': None,
            }