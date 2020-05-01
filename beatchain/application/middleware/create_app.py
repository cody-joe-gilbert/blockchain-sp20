# Blockchain & Applications Project 2: Beatchain
# Owner(s): Cody Gilbert

import os
import shutil
from typing import List
import asyncio
from hfc.fabric import Client
import middleware.constants as constants
from middleware.access_utils import enroll_user

async def create_channel(hf_client: Client,
                   install_org: str,
                   channel_name: str,
                   config_yaml_path: str,
                   ) -> Client:
    """
    Creates the given channel configuration artifacts and delivers them to the Orderer
    Args:
        hf_client: Network HF Client object
        install_org: Name of organization invoking channel creation
        channel_name: Name of the channel to create
        config_yaml_path: **Directory** in which the channel's 'configtx.yaml' file is stored. Does *not*
            take the path to the file itself.
    Returns:
        hf_client: Network HF Client object
    """
    # Create the HF client handle and extract config info
    network_info = hf_client.get_net_info()
    orderers = [x for x in network_info['orderers'].keys()]
    # Create a New Channel, the response should be true if succeed
    org_admin = hf_client.get_user(org_name=install_org, name='Admin')
    response = await hf_client.channel_create(
        orderer=orderers[0],
        channel_name=channel_name,
        requestor=org_admin,
        config_yaml=config_yaml_path,
        channel_profile=channel_name
    )
    if response:
        print(f'Created channel {channel_name}')
    else:
        raise ValueError(f'Failure to create channel {channel_name}')
    return hf_client

async def join_channel(hf_client: Client, channel_name: str) -> Client:
    """
    Joins each of the organization's peers to the newly created channel.

    Args:
        hf_client: Network HF Client object
        channel_name: Name of the channel to join
    Returns:
        hf_client: Network HF Client object
    """
    network_info = hf_client.get_net_info()
    orgs = [x for x in network_info['organizations'].keys()]
    orderers = [x for x in network_info['orderers'].keys()]

    # Each org joins the channel
    for org in orgs:
        if 'peers' not in network_info['organizations'][org]:
            # Skip the orderer
            continue
        peers = network_info['organizations'][org]['peers']
        # Join Peers into Channel, the response should be true if succeed
        org_admin = hf_client.get_user(org_name=org, name='Admin')
        responses = await hf_client.channel_join(
            requestor=org_admin,
            channel_name=channel_name,
            peers=peers,
            orderer=orderers[0]
        )
        if len(responses) == len(peers):
            print(f'Org {org} successfully joined channel {channel_name}')
        else:
            raise ValueError(f'Org {org} failed to joined channel {channel_name}')
    return hf_client

async def install_chaincode(hf_client: Client,
                      chaincode_store_path: str,
                      chaincode_gopath_rel_path: str,
                      chaincode_name: str,
                      chaincode_version: str) -> Client:
    """
    Joins each of the organization's peers to the newly created channel.
    Args:
        hf_client: Network HF Client object
        chaincode_store_path: Folder holding the chaincode.
            *Note:* the folder must have the structure <chaincode_store_path>/src/github.com/<actual chaincode dir>
        chaincode_gopath_rel_path: Path of the chaincode relative to the GOPATH. This is typically
            in the form 'github.com/<actual chaincode dir>'.
        chaincode_name: Name of the chaincode to install
        chaincode_version: Version number of the chaincode to install
    Returns:
        hf_client: Network HF Client object
    """

    gopath = os.path.normpath(os.path.join(
        os.path.dirname(os.path.realpath('__file__')),
        chaincode_store_path
    ))
    os.environ['GOPATH'] = os.path.abspath(gopath)
    network_info = hf_client.get_net_info()
    loop = asyncio.get_event_loop()
    orgs = [x for x in network_info['organizations'].keys()]
    # Install chaincode on each org's peers
    for org in orgs:
        if 'peers' not in network_info['organizations'][org]:
            # Skip the orderer
            continue
        peers = network_info['organizations'][org]['peers']
        # Join Peers into Channel, the response should be true if succeed
        org_admin = hf_client.get_user(org_name=org, name='Admin')
        responses = await hf_client.chaincode_install(
            requestor=org_admin,
            peers=peers,
            cc_path=chaincode_gopath_rel_path,
            cc_name=chaincode_name,
            cc_version=chaincode_version
        )
        if len(responses) == len(peers):
            print(f'Org {org} successfully installed chaincode')
        else:
            raise ValueError(f'Org {org} failed to install chaincode')
    return hf_client

async def instantiate_chaincode(hf_client,
                          install_org: str,
                          channel_name: str,
                          chaincode_name: str,
                          chaincode_version: str,
                          policy: dict,
                          instantiation_args: List[str] = []):
    """
    Instantiates the chaincode to create running containers on each of the peers.

    Args:
        hf_client: Network HF Client object
        install_org: Name of organization invoking chaincode instantiation
        channel_name: Name of the channel to on which to instantiate the chaincode
        chaincode_name: Name of the chaincode to instantiate
        chaincode_version: Version number of the chaincode to instantiate
        policy: A dictionary containing the endorsement policy as specified by HF. Note that this
            argument takes a dictionary directly, rather than a JSON-formatted string.
        instantiation_args: List of string arguments passed to the init method of the chaincode.
    Returns:
        hf_client: Network HF Client object
    """
    network_info = hf_client.get_net_info()
    org_admin = hf_client.get_user(org_name=install_org, name='Admin')
    response = await hf_client.chaincode_instantiate(
                                        requestor=org_admin,
                                        channel_name=channel_name,
                                        peers=[network_info['organizations'][install_org]['peers'][0]],
                                        args=instantiation_args,
                                        cc_name=chaincode_name,
                                        cc_version=chaincode_version,
                                        cc_endorsement_policy=policy,
                                        collections_config=None, # optional, for private data policy
                                        transient_map=None, # optional, for private data
                                        wait_for_event=True # optional, for being sure chaincode is instantiated
                                        )
    if response:
        print(f'Instantiated chaincode {chaincode_name} as org {install_org}')
    else:
        raise ValueError(f'Failure to instantiate chaincode {chaincode_name} as org {install_org}')

    return hf_client

async def create_app(install_org: str,
               admin_user_name: str,
               admin_password: str,
               test_mode: bool = False) -> None:
    """
    Bootstraps application creation by
    1. Creating the channel
    2. Joining all org peers to the channel
    3. Install the chaincode on all peers
    4. Instantiating the chaincode on all peers

    All network, channel, and chaincode settings are taken from constants
    stored in the middleware.constants module.

    Args:
        install_org: Name of organization creating the channel and chaincode instantiation
        admin_user_name: Username of the install_org admin
        admin_password: Password of the install_org admin
        test_mode: Flag to pass testing/debug arguments to the chaincode instantiation
    """
    if test_mode:
        args = constants.instantiation_args
    else:
        args = []

    # Cleanup any old existing key files stored in a file wallet
    if os.path.exists('./tmp/hfc-kvs'):
        shutil.rmtree('./tmp/hfc-kvs')

    # Create the HF client handle and extract config info
    hf_client = Client(net_profile=constants.config_path)

    # Validate a regular org admin status
    enroll_user(hf_client, install_org, admin_user_name, admin_password)

    hf_client = await create_channel(hf_client=hf_client,
                               install_org=install_org,
                               channel_name=constants.channel_name,
                               config_yaml_path=constants.config_yaml_path)

    hf_client = await join_channel(hf_client=hf_client, channel_name=constants.channel_name)

    hf_client = await install_chaincode(hf_client=hf_client,
                                  chaincode_store_path=constants.chaincode_store_path,
                                  chaincode_gopath_rel_path=constants.chaincode_gopath_rel_path,
                                  chaincode_name=constants.chaincode_name,
                                  chaincode_version=constants.chaincode_version)

    hf_client = await instantiate_chaincode(hf_client=hf_client,
                                      install_org=install_org,
                                      channel_name=constants.channel_name,
                                      chaincode_name=constants.chaincode_name,
                                      chaincode_version=constants.chaincode_version,
                                      policy=constants.policy,
                                      instantiation_args=args)

if __name__ == "__main__":
    asyncio.run(create_app(constants.install_org, 'admin', 'adminpw', True))




















