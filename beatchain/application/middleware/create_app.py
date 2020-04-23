import os
from typing import List
import asyncio
from hfc.fabric import Client
import middleware.constants as constants

def create_channel(hf_client,
                   install_org: str,
                   channel_name: str,
                   config_yaml_path: str,
                   ):
    loop = asyncio.get_event_loop()
    # Create the HF client handle and extract config info
    network_info = hf_client.get_net_info()
    orderers = [x for x in network_info['orderers'].keys()]


    # Create a New Channel, the response should be true if succeed
    org_admin = hf_client.get_user(org_name=install_org, name='Admin')
    response = loop.run_until_complete(hf_client.channel_create(
        orderer=orderers[0],
        channel_name=channel_name,
        requestor=org_admin,
        config_yaml=config_yaml_path,
        channel_profile=channel_name
    ))
    if response:
        print(f'Created channel {channel_name} as org {install_org}')
    else:
        raise ValueError(f'Failure to create channel {channel_name} as org {install_org}')
    return hf_client

def join_channel(hf_client, chaincode_store_path: str, channel_name: str):

    network_info = hf_client.get_net_info()
    loop = asyncio.get_event_loop()
    orgs = [x for x in network_info['organizations'].keys()]
    orderers = [x for x in network_info['orderers'].keys()]

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
        org_admin = hf_client.get_user(org_name=org, name='Admin')
        responses = loop.run_until_complete(hf_client.channel_join(
            requestor=org_admin,
            channel_name=channel_name,
            peers=peers,
            orderer=orderers[0]
        ))
        if len(responses) == len(peers):
            print(f'Org {org} successfully joined channel {channel_name}')
        else:
            raise ValueError(f'Org {org} failed to joined channel {channel_name}')
    return hf_client

def install_chaincode(hf_client,
                      chaincode_gopath_rel_path: str,
                      chaincode_name: str,
                      chaincode_version: str):

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
    return hf_client

def instantiate_chaincode(hf_client,
                          install_org: str,
                          channel_name: str,
                          chaincode_name: str,
                          chaincode_version: str,
                          policy: str,
                          instantiation_args: List[str] = []):
    # Instantiate the chaincode
    network_info = hf_client.get_net_info()
    loop = asyncio.get_event_loop()
    org_admin = hf_client.get_user(org_name=install_org, name='Admin')
    response = loop.run_until_complete(hf_client.chaincode_instantiate(
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
                                        ))
    if response:
        print(f'Instantiated chaincode {chaincode_name} as org {install_org}')
    else:
        raise ValueError(f'Failure to instantiate chaincode {chaincode_name} as org {install_org}')

    return hf_client

def create_app(install_org: str, test_mode: bool = False):

    if test_mode:
        args = constants.instantiation_args
    else:
        args = []

    # Create the HF client handle and extract config info
    hf_client = Client(net_profile=constants.config_path)

    hf_client = create_channel(hf_client=hf_client,
                               install_org=install_org,
                               channel_name=constants.channel_name,
                               config_yaml_path=constants.config_yaml_path)

    hf_client = join_channel(hf_client=hf_client,
                             chaincode_store_path=constants.chaincode_store_path,
                             channel_name=constants.channel_name)

    hf_client = install_chaincode(hf_client=hf_client,
                                  chaincode_gopath_rel_path=constants.chaincode_gopath_rel_path,
                                  chaincode_name=constants.chaincode_name,
                                  chaincode_version=constants.chaincode_version)

    hf_client = instantiate_chaincode(hf_client=hf_client,
                                      install_org=install_org,
                                      channel_name=constants.channel_name,
                                      chaincode_name=constants.chaincode_name,
                                      chaincode_version=constants.chaincode_version,
                                      policy=constants.policy,
                                      instantiation_args=args)

if __name__ == "__main__":
    create_app(constants.install_org, True)




















