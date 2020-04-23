install_org = 'appdevorg.beatchain.com'
config_path = '../../network/config.json'
config_yaml_path= '../../network/'  # Needs the folder for some reason
channel_name = 'fullchannel'
chaincode_name = 'beatchain'
chaincode_version = 'v0.1'
chaincode_store_path = '../../chaincode'
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


# Test Initialization args
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
    "init",
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