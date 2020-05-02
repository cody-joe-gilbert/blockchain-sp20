import asyncio
from middleware.access_utils import register_user
import middleware.constants as constants

loop = asyncio.get_event_loop()
user = constants.RegisterUserRequest(user_name='test4',
                                     admin_user_name='admin',
                                     admin_password='adminpw',
                                     role=None,
                                     affiliation=None,
                                     attrs=None)

secret = loop.run_until_complete(register_user(org_name='creatororg.beatchain.com', request=user))