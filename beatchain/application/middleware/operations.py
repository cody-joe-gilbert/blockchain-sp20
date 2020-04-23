import asyncio
import random
from typing import List
from hfc.fabric import Client
import middleware.constants as constants
from middleware.access_utils import enroll_user

def invoke(org_name: str,
           user_name: str,
           function: str,
           args: List[str]):

    # Invoke a chaincode
    hf_client = Client(net_profile=constants.config_path)
    network_info = hf_client.get_net_info()
    loop = asyncio.get_event_loop()
    peers = network_info['organizations'][org_name]['peers']
    random_peer = random.choice(peers)

    hf_client, user = enroll_user(hf_client, org_name, user_name)

    # The response should be true if succeed
    response = loop.run_until_complete(hf_client.chaincode_invoke(
        requestor=user,
        channel_name=constants.channel_name,
        peers=[random_peer],
        args=[function] + args,
        cc_name=constants.chaincode_name,
        transient_map=None, # optional, for private data
        wait_for_event=True, # for being sure chaincode invocation has been committed in the ledger, default is on tx event
    ))
    if not response:
        raise ValueError(f'Failure to invoke chaincode function {function}')
    return response


def query(org_name: str,
           user_name: str,
           function: str,
           args: List[str]):
    # Query a chaincode
    hf_client = Client(net_profile=constants.config_path)
    network_info = hf_client.get_net_info()
    loop = asyncio.get_event_loop()
    peers = network_info['organizations'][org_name]['peers']
    random_peer = random.choice(peers)

    hf_client, user = enroll_user(hf_client, org_name, user_name)

    # The response should be true if succeed
    response = loop.run_until_complete(hf_client.chaincode_query(
                                        requestor=user,
                                        channel_name=constants.channel_name,
                                        peers=[random_peer],
                                        args=[function] + args,
                                        cc_name=constants.chaincode_name
                                        ))
    if not response:
        raise ValueError(f'Failure to query chaincode function {function}')
    return response

