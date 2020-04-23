from middleware.create_app import create_app
from middleware.constants import constants

create_app(constants.install_org, True)