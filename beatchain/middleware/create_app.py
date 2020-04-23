import os
import asyncio
from hfc.fabric import Client
from hfc.fabric_ca.caservice import ca_service
from hfc.fabric_network import wallet

install_org = 'appdevorg.beatchain.com'
config_path = '../network/config.json'
config_yaml='../network/'  # Needs the folder for some reason
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

def enroll_admin(hf_client, org_name):
    network_info = hf_client.get_net_info()
    # Extract CA info
    org_info = network_info['organizations'][install_org]
    ca_name = org_info['certificateAuthorities'][0]
    ca_info = network_info['certificateAuthorities'][ca_name]

    casvc = ca_service(target=ca_info['url'])
    identity_service = casvc.newIdentityService()
    return casvc.enroll("admin", "adminpw")

def enroll_user(hf_client, org_name, user_name):
    network_info = hf_client.get_net_info()
    cred_wallet = wallet.FileSystenWallet()

    # Extract CA info
    org_info = network_info['organizations'][install_org]
    ca_name = org_info['certificateAuthorities'][0]
    ca_info = network_info['certificateAuthorities'][ca_name]

    # if user already exists, pull ID from storage
    if cred_wallet.exists(user_name):
        return create_user(user_name, org_name, org_info['mspid'])

    # If new user, then register and
    casvc = ca_service(target=ca_info['url'])
    identity_service = casvc.newIdentityService()
    admin_enrollment = casvc.enroll("admin", "adminpw")
    secret = admin_enrollment.register(user_name)
    user_enrollment = casvc.enroll(user_name, secret)

    # Store credentials in file kvs wallet; will be stored in ./tmp/hfc-kvs
    user_identity = wallet.Identity(user_name, user_enrollment)
    user_identity.CreateIdentity(cred_wallet)

    return cred_wallet.create_user(user_name, org_name, org_info['mspid'])


loop = asyncio.get_event_loop()
# Create the HF client handle and extract config info
cli = Client(net_profile=config_path)
network_info = cli.get_net_info()

orgs = [x for x in network_info['organizations'].keys()]
orderers = [x for x in network_info['orderers'].keys()]


# Create a New Channel, the response should be true if succeed
org_admin = cli.get_user(org_name=install_org, name='Admin')
response = loop.run_until_complete(cli.channel_create(
            orderer=orderers[0],
            channel_name=channel_name,
            requestor=org_admin,
            config_yaml=config_yaml,
            channel_profile=channel_name
            ))
if response:
    print(f'Created channel {channel_name} as org {install_org}')
else:
    raise ValueError(f'Failure to create channel {channel_name} as org {install_org}')


gopath_bak = os.environ.get('GOPATH', '')
gopath = os.path.normpath(os.path.join(
                      os.path.dirname(os.path.realpath('__file__')),
                      chaincode_store_path
                     ))
os.environ['GOPATH'] = os.path.abspath(gopath)


# Each org joins the channel
for org in orgs:
    if 'peers' not in network_info['organizations'][org]:
        # Skip the orderer
        continue
    peers = network_info['organizations'][org]['peers']
    # Join Peers into Channel, the response should be true if succeed
    org_admin = cli.get_user(org_name=org, name='Admin')
    responses = loop.run_until_complete(cli.channel_join(
                   requestor=org_admin,
                   channel_name=channel_name,
                   peers=peers,
                   orderer=orderers[0]
                   ))
    if len(responses) == len(peers):
        print(f'Org {org} successfully joined channel {channel_name}')
    else:
        raise ValueError(f'Org {org} failed to joined channel {channel_name}')

# Install chaincode on each org's peers
for org in orgs:
    if 'peers' not in network_info['organizations'][org]:
        # Skip the orderer
        continue
    peers = network_info['organizations'][org]['peers']
    # Join Peers into Channel, the response should be true if succeed
    org_admin = cli.get_user(org_name=org, name='Admin')
    responses = loop.run_until_complete(cli.chaincode_install(
                   requestor=org_admin,
                   peers=peers,
                   cc_path=chaincode_gopath_rel_path,
                   cc_name=chaincode_name,
                   cc_version=chaincode_version
                   ))
    if len(responses) == len(peers):
        print(f'Org {org} successfully installed chaincode')
    else:
        raise ValueError(f'Org {org} failed to install chaincode')

# Instantiate the chaincode
org_admin = cli.get_user(org_name=install_org, name='Admin')
response = loop.run_until_complete(cli.chaincode_instantiate(
               requestor=org_admin,
               channel_name=channel_name,
               peers=[network_info['organizations'][install_org]['peers'][0]],
               args=[],
               cc_name=chaincode_name,
               cc_version=chaincode_version,
               cc_endorsement_policy=policy,
               collections_config=None, # optional, for private data policy
               transient_map=None, # optional, for private data
               wait_for_event=True # optional, for being sure chaincode is instantiated
               ))
if response:
    print(f'Instantiated chaincode {chaincode_name} as org {install_org}')
else:
    raise ValueError(f'Failure to instantiate chaincode {chaincode_name} as org {install_org}')







