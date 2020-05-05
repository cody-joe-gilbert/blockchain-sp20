import asyncio
import random
from hfc.fabric import Client
import middleware.constants as constants
from middleware.access_utils import enroll_user


org_name = 'appdevorg.beatchain.com'
user_name = 'admin'
user_password = 'adminpw'
function = 'ListBankAccounts'


hf_client = Client(net_profile=constants.config_path)

hf_client.new_channel(constants.channel_name)

network_info = hf_client.get_net_info()

peers = network_info['organizations'][org_name]['peers']
random_peer = random.choice(peers)

user = enroll_user(hf_client, org_name, user_name, user_password)

loop = asyncio.get_event_loop()
response = loop.run_until_complete(hf_client.chaincode_query(requestor=user,
                                           channel_name=constants.channel_name,
                                           fcn=function,
                                           peers=[random_peer],
                                           args=[],
                                           cc_name=constants.chaincode_name))

print(response)