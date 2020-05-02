# Blockchain & Applications Project 2: Beatchain
# Owner(s): Cody Gilbert

import os
import hfc.fabric.user
import middleware.constants as constants
from hfc.fabric import Client
from hfc.fabric_ca.caservice import ca_service
from hfc.fabric_network import wallet

async def register_user(org_name: str, request: constants.RegisterUserRequest) -> str:
    """
    Registers a user to the Org's Fabric CA Server
    Args:
        org_name: Organization's name
        request: RegisterUserRequest object containing
            registration information
    Returns:
        Pre-generated user secret
    """
    # Create/Open a wallet on a temp path including the org name
    # Org name must be included, otherwise usernames must be unique
    # over all orgs
    wallet_path = os.path.join(os.getcwd(), 'tmp', 'hfc-kvs', org_name)
    cred_wallet = wallet.FileSystenWallet(path=wallet_path)  # [sic]

    # Setup a HF network client
    hf_client = Client(net_profile=constants.config_path)
    hf_client.new_channel(constants.channel_name)

    # Extract CA info
    network_info = hf_client.get_net_info()
    org_info = network_info['organizations'][org_name]
    ca_name = org_info['certificateAuthorities'][0]
    ca_info = network_info['certificateAuthorities'][ca_name]

    # if user already exists, pull ID from storage
    if cred_wallet.exists(request.user_name):
        return None
    casvc = ca_service(target=ca_info['url'])
    admin_enrollment = casvc.enroll(request.admin_user_name, request.admin_password)

    secret = admin_enrollment.register(enrollmentID=request.user_name,
                                       enrollmentSecret=request.user_password,
                                       role=request.role,
                                       affiliation=request.affiliation,
                                       attrs=request.attrs)

    return secret

def enroll_user(hf_client: hfc.fabric.Client,
                org_name: str,
                user_name: str,
                user_password: str,
                ) -> hfc.fabric.user.User:
    """
    Enrolls a user to the Org's Fabric CA Server
    Args:
        hf_client: Network HF Client object
        org_name: Organization's name
        user_name: Username to enroll
        user_password: User's password
    Returns:
        Enrolled User object
    """
    # Create/Open a wallet on a temp path including the org name
    # Org name must be included, otherwise usernames must be unique
    # over all orgs
    wallet_path = os.path.join(os.getcwd(), 'tmp', 'hfc-kvs', org_name)
    cred_wallet = wallet.FileSystenWallet(path=wallet_path)  # [sic]

    # Extract CA info
    network_info = hf_client.get_net_info()
    org_info = network_info['organizations'][org_name]
    ca_name = org_info['certificateAuthorities'][0]
    ca_info = network_info['certificateAuthorities'][ca_name]

    # if user already exists, pull ID from storage
    if cred_wallet.exists(user_name):
        user = cred_wallet.create_user(user_name, org_name, org_info['mspid'])
        #if user.enrollment_secret != user_password:
        #    # TODO: Check passwords in a *much* more secure way than this
        #    raise ValidationError('Invalid username/password')
        return user

    casvc = ca_service(target=ca_info['url'])
    user_enrollment = casvc.enroll(user_name, user_password)

    # Store credentials in file kvs wallet; will be stored in ./tmp/hfc-kvs
    user_identity = wallet.Identity(user_name, user_enrollment)
    user_identity.CreateIdentity(cred_wallet)

    return cred_wallet.create_user(user_name, org_name, org_info['mspid'])