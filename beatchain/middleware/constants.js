
var os = require('os');
var path = require('path');

var tempdir = "../network/client-certs";
//path.join(os.tmpdir(), 'hfc');

// Frame the endorsement policy
var FULL_CHANNEL_POLICY = [{
	role: {
		name: 'member',
		mspId: 'AppDevMSP'
	}
}, {
	role: {
		name: 'member',
		mspId: 'CreatorMSP'
	}
}, {
	role: {
		name: 'member',
		mspId: 'CustomerMSP'
	}
}, {
	role: {
		name: 'admin',
		mspId: 'BeatchainMSP'
	}
}];

var ONE_OF_THREE_ORG_MEMBER = {
	identities: FULL_CHANNEL_POLICY,
	policy: {
		'1-of': [{ 'signed-by': 0 }, { 'signed-by': 1 }, { 'signed-by': 2 }]
	}
};

var ALL_THREE_ORG_MEMBERS = {
	identities: FULL_CHANNEL_POLICY,
	policy: {
		'3-of': [{ 'signed-by': 0 }, { 'signed-by': 1 }, { 'signed-by': 2 }]
	}
};

var ACCEPT_ALL = {
	identities: [],
	policy: {
		'0-of': []
	}
};

var chaincodeLocation = '../chaincode';

var networkId = 'beatchain';

var networkConfig = './config.json';

var networkLocation = '../network';

var channelConfig = 'channel-artifacts/channel.tx';

var APPDEV_ORG = 'appdevorg';
var CREATOR_ORG = 'creatororg';
var CUSTOMER_ORG = 'customerorg';

var CHANNEL_NAME = 'fullchannel';
var CHAINCODE_PATH = '';
var CHAINCODE_ID = 'beatchain_alpha';
var CHAINCODE_VERSION = 'v0';

var TRANSACTION_ENDORSEMENT_POLICY = ALL_THREE_ORG_MEMBERS;

module.exports = {
	tempdir: tempdir,
	chaincodeLocation: chaincodeLocation,
	networkId: networkId,
	networkConfig: networkConfig,
	networkLocation: networkLocation,
	channelConfig: channelConfig,
	APPDEV_ORG: APPDEV_ORG,
	CREATOR_ORG: CREATOR_ORG,
	CUSTOMER_ORG: CUSTOMER_ORG,
	CHANNEL_NAME: CHANNEL_NAME,
	CHAINCODE_PATH: CHAINCODE_PATH,
	CHAINCODE_ID: CHAINCODE_ID,
	CHAINCODE_VERSION: CHAINCODE_VERSION,
	ALL_THREE_ORG_MEMBERS: ALL_THREE_ORG_MEMBERS,
	TRANSACTION_ENDORSEMENT_POLICY: TRANSACTION_ENDORSEMENT_POLICY
};
