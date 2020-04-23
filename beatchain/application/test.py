from middleware.create_app import create_app
import middleware.constants as constants

create_app(constants.install_org, True)