
import asyncio
import middleware.constants as constants
import middleware.access_utils as access_utils
import middleware.operations as operations

admin = {'username': 'admin',
         'password': 'adminpw'}

creator_member = {'username': 'testCreator',
                  'password': 'creatorpw',
                  'id': "3333",
                  'org': 'creatororg.beatchain.com'
                  }
customer_member = {'username': 'testCustomer',
                  'password': 'customerpw',
                  'id': "2222",
                  'org': 'customerorg.beatchain.com'
                  }
appdev_member = {'username': 'testAppdev',
                 'password': 'appdevpw',
                 'id': "1111",
                 'org': 'appdevorg.beatchain.com'}

loop = asyncio.get_event_loop()


print('Registering Creator user')
register_req = constants.RegisterUserRequest(
    admin_user_name=admin['username'],
    admin_password=admin['password'],
    user_name=creator_member['username'],
    user_password=creator_member['password'],
    role='client',
    attrs=[{'name':'id', 'value': creator_member['id']}])
loop.run_until_complete(access_utils.register_user(creator_member['org'], register_req))
print('Creator user registered')

print('Registering Customer user')
register_req = constants.RegisterUserRequest(
    admin_user_name=admin['username'],
    admin_password=admin['password'],
    user_name=customer_member['username'],
    user_password=customer_member['password'],
    role='client',
    attrs=[{'name':'id', 'value': customer_member['id']}])
loop.run_until_complete(access_utils.register_user(customer_member['org'], register_req))
print('Customer user registered')


print('Registering AppDev user')
register_req = constants.RegisterUserRequest(
    admin_user_name=admin['username'],
    admin_password=admin['password'],
    user_name=appdev_member['username'],
    user_password=appdev_member['password'],
    role='client',
    attrs=[{'name':'id', 'value': appdev_member['id']}])
loop.run_until_complete(access_utils.register_user(appdev_member['org'], register_req))
print('AppDev user registered')


print('Creating a new product')
product_id = loop.run_until_complete(operations.invoke(creator_member['org'],
                                                       creator_member['username'],
                                                       creator_member['password'],
                                                       constants.channel_name,
                                                       function='AddProduct',
                                                       args=['Test Product Name']))
print('New Product created with ID: ', product_id)


print('Creating a new contract')
loop.run_until_complete(operations.invoke(appdev_member['org'],
                                                       appdev_member['username'],
                                                       appdev_member['password'],
                                                       constants.channel_name,
                                                       function='OfferContract',
                                                       args=[appdev_member['id'],
                                                             creator_member['id'],
                                                             product_id,
                                                             '0.02']))
print('New Contract offered')


print('Accept Contract')
loop.run_until_complete(operations.invoke(creator_member['org'],
                                                       creator_member['username'],
                                                       creator_member['password'],
                                                       constants.channel_name,
                                                       function='AcceptContract',
                                                       args=[creator_member['id'],
                                                             product_id,
                                                             appdev_member['id']]))
print('Contract Accepted')




print('Request Song')
loop.run_until_complete(operations.invoke(customer_member['org'],
                                                       customer_member['username'],
                                                       customer_member['password'],
                                                       constants.channel_name,
                                                       function='RequestSong',
                                                       args=[product_id]))
print('Song Requested')



print('Creator Collect Payment')
loop.run_until_complete(operations.invoke(creator_member['org'],
                                                       creator_member['username'],
                                                       creator_member['password'],
                                                       constants.channel_name,
                                                       function='CollectPayment',
                                                       args=[]))
print('Payment Collected')