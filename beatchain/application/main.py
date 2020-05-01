# Blockchain & Applications Project 2: Beatchain
# Owner(s): Cody Gilbert

from typing import List
from pydantic import BaseModel, Field
from fastapi import FastAPI, Query
import middleware.constants as constants
import middleware.access_utils as access_utils
from middleware.create_app import create_app
import middleware.operations as operations

app = FastAPI()

# These request classes are used by FastAPI to validate and parse
# the request body to the API endpoint
class CreateAppRequest(BaseModel):
    admin_user_name: str
    admin_password: str

@app.post('/admin/create_app')
def creation_request(req: CreateAppRequest,
                     org_name: constants.OrgNames = Query(..., title="Organization Name"),
                     test_mode: bool = Query(False, title="Debug Initialization Mode Flag"),
                     ):
    """
    Submits a request to bootstrap the application on the
    HF network.
    """
    try:
        create_app(org_name,
                   req.admin_user_name,
                   req.admin_password,
                   test_mode)
    except Exception as e:
        return {'Status': 'Application Creation Request failed',
                'Error': repr(e)}
    return {'Status': 'Application Created',
            'Error': None}

class RegisterUserRequest(BaseModel):
    user_name: str
    admin_user_name: str
    admin_password: str

@app.post('/admin/register')
def register(req: RegisterUserRequest,
             org_name: constants.OrgNames = Query(..., title="Organization Name"),
             ):
    """
    Submits a request to register a user on HF network using the given
    org's certificate authority. If successful, return's the user's
    login password with the response.
    """
    # TODO: Passing a secret back is NOT secure! This section is for demo only!
    try:
        secret = access_utils.register_user(org_name,
                                            req.user_name,
                                            req.admin_user_name,
                                            req.admin_password)
    except Exception as e:
        return {'Status': 'Registration Request failed',
                'Secret': None,
                'Error': repr(e)}
    return {'Status': 'Application Created',
            'Secret': secret,
            'Error': None}

class InvokeRequest(BaseModel):
    user_name: str
    user_password: str
    args: List[str] = []

@app.post('/invoke')
async def invoke_request(req: InvokeRequest,
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

class QueryRequest(BaseModel):
    user_name: str
    user_password: str
    args: List[str] = []

@app.post('/query')
async def query_request(req: QueryRequest,
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
