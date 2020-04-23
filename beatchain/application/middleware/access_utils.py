from hfc.fabric_ca.caservice import ca_service
from hfc.fabric_network import wallet


def enroll_admin(hf_client, org_name):
    network_info = hf_client.get_net_info()
    # Extract CA info
    org_info = network_info['organizations'][org_name]
    ca_name = org_info['certificateAuthorities'][0]
    ca_info = network_info['certificateAuthorities'][ca_name]

    casvc = ca_service(target=ca_info['url'])
    identity_service = casvc.newIdentityService()
    return hf_client, casvc.enroll("admin", "adminpw")

def enroll_user(hf_client, org_name, user_name):
    network_info = hf_client.get_net_info()
    cred_wallet = wallet.FileSystenWallet()

    # Extract CA info
    org_info = network_info['organizations'][org_name]
    ca_name = org_info['certificateAuthorities'][0]
    ca_info = network_info['certificateAuthorities'][ca_name]

    # if user already exists, pull ID from storage
    if cred_wallet.exists(user_name):
        return cred_wallet.create_user(user_name, org_name, org_info['mspid'])

    # If new user, then register and
    casvc = ca_service(target=ca_info['url'])
    hf_client, admin_enrollment = enroll_admin(hf_client, org_name)
    secret = admin_enrollment.register(user_name)
    user_enrollment = casvc.enroll(user_name, secret)

    # Store credentials in file kvs wallet; will be stored in ./tmp/hfc-kvs
    user_identity = wallet.Identity(user_name, user_enrollment)
    user_identity.CreateIdentity(cred_wallet)

    return hf_client, cred_wallet.create_user(user_name, org_name, org_info['mspid'])