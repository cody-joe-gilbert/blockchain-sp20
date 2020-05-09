import asyncio
import random
from hfc.fabric import Client
import middleware.constants as constants
import middleware.operations as operations
from middleware.access_utils import enroll_user


org_name = 'appdevorg.beatchain.com'
admin_user_name = 'admin'
admin_password = 'adminpw'
function = 'ListBankAccounts'
admin_fee_frac = 0.5

hf_client = Client(net_profile=constants.config_path)

hf_client.new_channel(constants.channel_name)

network_info = hf_client.get_net_info()

peers = network_info['organizations'][org_name]['peers']
random_peer = random.choice(peers)

user = enroll_user(hf_client, org_name, user_name, user_password)
user = hf_client.get_user(org_name=org_name, name='Admin')

loop = asyncio.get_event_loop()

response = loop.run_until_complete(operations.invoke('appdevorg.beatchain.com',
                                                     admin_user_name,
                                                     admin_password,
                                                     constants.channel_name,
                                                     function='AddAppDevRecord',
                                                     args=["afldsfj;ls"]))



# Query Peer installed chaincodes, make sure the chaincode is installed
response = loop.run_until_complete(operations.invoke('appdevorg.beatchain.com',
                                                     admin_user_name,
                                                     admin_password,
                                                     constants.channel_name,
                                                     function='AddAppDevRecord',
                                                     args=[str(round(admin_fee_frac, 3))]))





response = loop.run_until_complete(hf_client.query_info(
    requestor=user,
    channel_name='dud',
    peers=[random_peer]
))

response = loop.run_until_complete(hf_client. query_installed_chaincodes(
requestor=user,
peers=[random_peer]
))

print(response)