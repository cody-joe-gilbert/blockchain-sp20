# Blockchain & Applications Project 2: Beatchain
# Owner(s): Cody Gilbert

import hfc.fabric.user
import middleware.constants as constants
from hfc.fabric import Client
from hfc.fabric_ca.caservice import ca_service
from hfc.fabric_network import wallet

async def register_user(org_name: str,
                  user_name: str,
                  admin_user_name: str,
                  admin_password: str) -> str:
    """
    Registers a user to the Org's Fabric CA Server
    Args:
        org_name: Organization's name
        user_name: Username to register
        admin_user_name: Org admin username
        admin_password: Org admin password
    Returns:
        Pre-generated user secret
    """
    # Setup a HF network client
    hf_client = Client(net_profile=constants.config_path)

    network_info = hf_client.get_net_info()
    cred_wallet = wallet.FileSystenWallet()

    # Extract CA info
    org_info = network_info['organizations'][org_name]
    ca_name = org_info['certificateAuthorities'][0]
    ca_info = network_info['certificateAuthorities'][ca_name]

    # if user already exists, pull ID from storage
    if cred_wallet.exists(user_name):
        return None
    casvc = ca_service(target=ca_info['url'])
    admin_enrollment = casvc.enroll(admin_user_name, admin_password)
    return admin_enrollment.register(user_name)

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

    cred_wallet = wallet.FileSystenWallet()

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