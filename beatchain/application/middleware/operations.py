# Blockchain & Applications Project 2: Beatchain
# Owner(s): Cody Gilbert

import random
from typing import List, Optional, Dict
from hfc.fabric import Client
import middleware.constants as constants
from middleware.access_utils import enroll_user

async def invoke(org_name: str,
           user_name: str,
           user_password: str,
           channel_name: str,
           function: str,
           args: List[str]) -> str:
    """
    Submits a blockchain transaction invocation to all the peers in
    the network.

    Args:
        org_name: Name of the submitting user's organization
        user_name: Username of the submitter
        user_password: Password of the submitting user
        channel_name: Name of the channel on which to connect client
        function: Name of the chaincode function to invoke
        args: A list of string arguments passed to the chaincode
    Returns:
        Response string from the *first* peer that responded with
        confirmation of the execution.
    """

    # Setup a HF network client
    hf_client = Client(net_profile=constants.config_path)

    # Connect to the given channel
    hf_client.new_channel(channel_name)

    # Gather information about the network
    network_info = hf_client.get_net_info()

    # Invocations require read/write sets from as many peers
    # as specified in the endorsement policy. Here we will
    # go ahead and request endorsement from ALL peers.
    # If your policy requires only a subset, you may wish to
    # alter this section
    peers = []
    for org in network_info['organizations'].keys():
        if 'peers' in network_info['organizations'][org]:
            for peer in network_info['organizations'][org]['peers']:
                peers.append(peer)

    # Enroll the user that will be invoking the query
    user = enroll_user(hf_client, org_name, user_name, user_password)

    # Submit the query to the peers and await a response
    response = await hf_client.chaincode_invoke(requestor=user,
                                                channel_name=channel_name,
                                                peers=peers,
                                                fcn=function,
                                                args=args,
                                                cc_name=constants.chaincode_name,
                                                transient_map=None, # optional, for private data
                                                wait_for_event=True, # for being sure chaincode invocation has been commited in the ledger, default is on tx event
                                               )

    if not response:
        raise ValueError(f'Failure to invoke chaincode function {function} with response: {response}')
    return response


async def query(org_name: str,
          user_name: str,
          user_password: str,
          channel_name: str,
          function: str,
          args: List[str]) -> str:
    """
    Submits a ledger query to a single peer within the specified org.
    Note that queries will NOT submit any changes to the ledger state,
    even if an invocation query function is invoked.
    Args:
        org_name: Name of the submitting user's organization
        user_name: Username of the submitter
        user_password: Password of the submitting user
        channel_name: Name of the channel on which to connect client
        function: Name of the chaincode function to invoke
        args: A list of string arguments passed to the chaincode
    Returns:
        Response string from the given peer
    """
    # Setup a HF network client
    hf_client = Client(net_profile=constants.config_path)

    # Connect to the given channel
    hf_client.new_channel(channel_name)

    # Gather information about the network
    network_info = hf_client.get_net_info()

    # For queries, we only need a single peer
    # Here we will randomly select one from our org
    peers = network_info['organizations'][org_name]['peers']
    random_peer = random.choice(peers)

    # Enroll the user that will be invoking the query
    user = enroll_user(hf_client, org_name, user_name, user_password)

    # Submit the query to the selected peer and await a response
    response = await hf_client.chaincode_query(requestor=user,
                                                channel_name=channel_name,
                                                fcn=function,
                                                peers=[random_peer],
                                                args=args,
                                                cc_name=constants.chaincode_name)
    if not response:
        raise ValueError(f'Failure to query chaincode function {function} with response: {response}')
    return response

def get_network_info() -> Dict:
    '''
    Returns the dictionary containing the network information
    Returns:
        Dictionary of network configuration information
    '''
    hf_client = Client(net_profile=constants.config_path)
    return hf_client.get_net_info()

async def get_channels(org_name: str,
                 user_name: str,
                 user_password: str) -> List[str]:
    """
    Fetches the IDs of all the channels on the network
    Args:
        org_name: Name of the submitting user's organization
        user_name: Username of the submitter
        user_password: Password of the submitting user
    Returns:
        List of channel IDs
    """
    # Setup a HF network client
    hf_client = Client(net_profile=constants.config_path)

    # Gather information about the network
    network_info = hf_client.get_net_info()

    # For queries, we only need a single peer
    # Here we will randomly select one from our org
    peers = network_info['organizations'][org_name]['peers']
    random_peer = random.choice(peers)

    user = enroll_user(hf_client, org_name, user_name, user_password)

    response = await hf_client.query_channels(
        requestor=user,
        peers=[random_peer]
    )
    return [x.channel_id for x in response.channels]

async def get_block_info(org_name: str,
                       user_name: str,
                       user_password: str,
                       channel_name: str = constants.channel_name) -> str:
    """
    Fetches the blockchain state of the given channel
    Args:
        org_name: Name of the submitting user's organization
        user_name: Username of the submitter
        user_password: Password of the submitting user
        channel_name: Name of the channel on which to connect client
    Returns:
        String of blockchain info
    """
    # Setup a HF network client
    hf_client = Client(net_profile=constants.config_path)

    # Connect to the given channel
    hf_client.new_channel(channel_name)

    # Gather information about the network
    network_info = hf_client.get_net_info()

    # For queries, we only need a single peer
    # Here we will randomly select one from our org
    peers = network_info['organizations'][org_name]['peers']
    random_peer = random.choice(peers)

    user = enroll_user(hf_client, org_name, user_name, user_password)

    response = await hf_client.query_info(
        requestor=user,
        channel_name=channel_name,
        peers=[random_peer]
    )
    return str(response)

async def get_instantiated_chaincodes(org_name: str,
                         user_name: str,
                         user_password: str,
                         channel_name: str = constants.channel_name) -> List[Dict[str, str]]:
    """
    Fetches the blockchain state of the given channel
    Args:
        org_name: Name of the submitting user's organization
        user_name: Username of the submitter
        user_password: Password of the submitting user
        channel_name: Name of the channel on which to connect client
    Returns:
        List of instantiated chaincodes, each of which contains:
            'Name': Chaincode name,
            'Version': Chaincode version,
            'Path': GOPATH relative chaincode location
    """
    # Setup a HF network client
    hf_client = Client(net_profile=constants.config_path)

    # Connect to the given channel
    hf_client.new_channel(channel_name)

    # Gather information about the network
    network_info = hf_client.get_net_info()

    # For queries, we only need a single peer
    # Here we will randomly select one from our org
    peers = network_info['organizations'][org_name]['peers']
    random_peer = random.choice(peers)

    user = enroll_user(hf_client, org_name, user_name, user_password)

    response = await hf_client.query_instantiated_chaincodes(
        requestor=user,
        channel_name=channel_name,
        peers=[random_peer]
    )
    chaincodes = []
    for chaincode in response[0].chaincodes:
        code = {'Name': chaincode.name,
                'Version': chaincode.version,
                'Path': chaincode.path}
        chaincodes.append(code)

    return chaincodes
