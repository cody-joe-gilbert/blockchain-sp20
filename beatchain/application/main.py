from fastapi import FastAPI
from middleware.create_app import create_app
from middleware.operations import (invoke, query)

app = FastAPI()

@app.get('/create_app')
async def creation_request(org_name: str, test_mode: bool = False):
    try:
        create_app(org_name, test_mode)
    except Exception as e:
        return {'Status': 'Creation Request failed',
                'Error': repr(e)}
    return {'Status': 'Application Created',
            'Error': None}

@app.get('/invoke')
async def invoke_request(function: str, org_name: str, user_name: str, args: List[str]):
    try:
        response = invoke(org_name, user_name, function, args)
    except Exception as e:
        return {'Status': 'Invoke Request failed',
                'Response': None,
                'Error': repr(e),
                }
    return {'Status': 'Invoke Request successful',
            'Response': response,
            'Error': None,
            }

@app.get('/query')
async def query_request(function: str, org_name: str, user_name: str, args: List[str]):
    try:
        response = query(org_name, user_name, function, args)
    except Exception as e:
        return {'Status': 'Query Request failed',
                'Response': None,
                'Error': repr(e),
                }
    return {'Status': 'Query Request failed',
            'Response': response,
            'Error': None,
            }