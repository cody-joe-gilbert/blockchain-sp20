# Blockchain & Applications Project 2: Beatchain
# Owner(s): Cody Gilbert
from typing import Optional, Dict, List
from enum import Enum
from pydantic import BaseModel

# Pydantic Enum Validation Classes
class InvokeFunctions(str, Enum):
    """
    Chaincode functions supporting 'invoke' transactions
    Note: These functions must match those specified in the
    chaincode 'invoke' function
    """
    RenewSubscription = "RenewSubscription"
    CollectPayment = "CollectPayment"
    OfferContract = "OfferContract"
    AcceptContract = "AcceptContract"
    RejectContract = "RejectContract"

class QueryFunctions(str, Enum):
    """
    Chaincode functions supporting 'query' transactions
    Note: These functions must match those specified in the
    chaincode 'invoke' function
    """
    ListBankAccounts = "ListBankAccounts"
    ListCustomers = "ListCustomers"

class OrgNames(str, Enum):
    """
    Organization Names
    Note: These names must match those specified in the
    configuration files within ../network
    """
    beatchainorg = "beatchainorg.beatchain.com"
    appdevorg = "appdevorg.beatchain.com"
    creatororg = "creatororg.beatchain.com"
    customerorg = "customerorg.beatchain.com"

class ChannelNames(str, Enum):
    """
    HF Network Channel names
    Note: These names must match those specified in the
    configuration files within ../network
    """
    fullchannel = "fullchannel"

class UserRoles(str, Enum):
    """
    HF user roles
    Note: These names must match those specified in the
    Fabric CA server configuration file fabric-ca-server-config.yaml
    """
    admin = "admin"
    client = "client"

class UserAffiliations(str, Enum):
    """
    HF user roles
    Note: These names must match those specified in the
    Fabric CA server configuration file fabric-ca-server-config.yaml
    """
    # beatchain = "beatchain"
    # customer = "customer"
    # creator = "creator"
    # dotify = "dotify"
    # waval = "waval"
    # cantcloseboxa = "cantcloseboxa"
    # sometimesytunes = "sometimesytunes"
    org1_department1 = "org1.department1"

# These request classes are used by FastAPI to validate and parse
# the request body to the API endpoint

class UserRegistrationAttrs(BaseModel):
    name: str
    value: str

class RegisterUserRequest(BaseModel):
    admin_user_name: str
    admin_password: str
    user_name: str
    user_password: Optional[str]
    role: Optional[UserRoles]
    affiliation: Optional[UserAffiliations]
    attrs: Optional[List[UserRegistrationAttrs]]

class CreateAppRequest(BaseModel):
    admin_user_name: str
    admin_password: str

class InvokeRequest(BaseModel):
    user_name: str
    user_password: str
    args: List[str] = []

class AddProductRequest(BaseModel):
    user_name: str
    user_password: str
    product_name: str

class AddUserRecordRequest(BaseModel):
    admin_user_name: str
    admin_password: str
    user_name: str
    user_password: Optional[str]
    affiliation: Optional[UserAffiliations]

class AddCustomerRecordRequest(BaseModel):
    admin_user_name: str
    admin_password: str
    user_name: str
    user_password: Optional[str]
    affiliation: Optional[UserAffiliations]



# App creation parameters
install_org = 'appdevorg.beatchain.com'
config_path = '../network/config.json'
config_yaml_path= '../network/'  # Needs the folder for some reason
channel_name = 'fullchannel'
chaincode_name = 'beatchain'
chaincode_version = 'v0.1'
chaincode_store_path = '../chaincode'
chaincode_gopath_rel_path = 'github.com/beatchain'

FULL_CHANNEL_POLICY = [{
    'role': {
        'name': 'member',
        'mspId': 'AppDevMSP'
    }
}, {
    'role': {
        'name': 'member',
        'mspId': 'CreatorMSP'
    }
}, {
    'role': {
        'name': 'member',
        'mspId': 'CustomerMSP'
    }
}, {
    'role': {
        'name': 'admin',
        'mspId': 'BeatchainMSP'
    }
}]

policy = { 'identities': FULL_CHANNEL_POLICY,
           'policy': {
               '3-of': [{ 'signed-by': 0 }, { 'signed-by': 1 }, { 'signed-by': 2 }]
           }
           }

# Test Initialization input arguments
BEATCHAIN_ADMIN_BALANCE = "1000"
TEST_APPDEV_ID = "1111"
TEST_APPDEV_BA_ID = "1111"
TEST_APPDEV_DEVSHARE = "0.1"
TEST_APPDEV_BA_BALANCE = "1000"
TEST_CUSTOMER_ID = "2222"
TEST_CUSTOMER_BA_ID = "2222"
TEST_CUSTOMER_SUBFEE = "1.00"
TEST_CUSTOMER_SUB_DUE_DATE = "2020-06-01"
TEST_CUSTOMER_BA_BALANCE = "1000"
TEST_CREATOR_ID = "3333"
TEST_CREATOR_BA_ID = "3333"
TEST_CREATOR_BA_BALANCE = "1000"
TEST_PRODUCT_ID = "4444"
TEST_PRODUCT_NAME = "Test Product"
TEST_PRODUCT_TOTAL_LISTENS = "5"
TEST_PRODUCT_UNREN_LISTENS = "3"
TEST_PRODUCT_TOTAL_METRICS = "7"
TEST_PRODUCT_UNREN_METRICS = "4"
TEST_PRODUCT_ADD_METRICS = "0"
TEST_PRODUCT_ACTIVE = "true"
TEST_CONTRACT_PPS = "0.01"
TEST_CONTRACT_STATUS = "true"

instantiation_args = [
    BEATCHAIN_ADMIN_BALANCE,
    TEST_APPDEV_ID,
    TEST_APPDEV_BA_ID,
    TEST_APPDEV_DEVSHARE,
    TEST_APPDEV_BA_BALANCE,
    TEST_CUSTOMER_ID,
    TEST_CUSTOMER_BA_ID,
    TEST_CUSTOMER_SUBFEE,
    TEST_CUSTOMER_SUB_DUE_DATE,
    TEST_CUSTOMER_BA_BALANCE,
    TEST_CREATOR_ID,
    TEST_CREATOR_BA_ID,
    TEST_CREATOR_BA_BALANCE,
    TEST_PRODUCT_ID,
    TEST_PRODUCT_NAME,
    TEST_PRODUCT_TOTAL_LISTENS,
    TEST_PRODUCT_UNREN_LISTENS,
    TEST_PRODUCT_TOTAL_METRICS,
    TEST_PRODUCT_UNREN_METRICS,
    TEST_PRODUCT_ADD_METRICS,
    TEST_PRODUCT_ACTIVE,
    TEST_CONTRACT_PPS,
    TEST_CONTRACT_STATUS,
]