# Blockchain & Applications Project 2: Beatchain
# Owner(s): Cody Gilbert
# Note: Method here uses Python process spawning to execute Fabric JS SDK scripts. This would not be
#   acceptable in a production environment, but is considered acceptable here as it can be replaced with
#   the Fabric Python SDK functions once those tools are released as stable.

from fastapi import FastAPI
from Naked.toolshed.shell import muterun_js

app = FastAPI()

class OrgName(str, Enum):
    beatchainorg = "beatchainorg"
    appdevorg = "appdevorg"
    creatororg = "creatororg"
    customerorg = "customerorg"

@app.get('/query/ListBankAccounts')
def query_ListBankAccounts(org_name: OrgName, user_name: str = 'admin'):
    '''
    Fetches all bank accounts records from Beatchain.
    Args:
        org_name: Name of the organization accessing the records
        user_name: Login ID of the user accessing the records

    '''
    js_script = '../middleware/query_wrapper.js'
    function = 'ListBankAccounts'
    js_args = ' '.join([org_name, function, '\"\"', user_name])
    response = muterun_js(js_script, js_args)
    if response.exitcode == 0:
        return {'Status': 'Query Request Successful',
                'Response': response.stdout,
                'Error': None,
                }
    else:
        return {'Status': 'Query Request Failed',
                'Response': response.stdout,
                'Error': response.stderr,
                }
