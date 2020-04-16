/*
 * Copyright 2018 IBM All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the 'License');
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an 'AS IS' BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

'use strict';

var Constants = require('./constants.js');
var ClientUtils = require('./clientUtils.js');
var createChannel = require('./create-channel.js');
var joinChannel = require('./join-channel.js');
var installCC = require('./install-chaincode.js');
var instantiateCC = require('./instantiate-chaincode.js');
var invokeCC = require('./invoke-chaincode.js');
var queryCC = require('./query-chaincode.js');

// Create a channel using the given network configuration
createChannel.createChannel(Constants.CHANNEL_NAME).then(() => {
	console.log('\n');
	console.log('--------------------------');
	console.log('CHANNEL CREATION COMPLETE');
	console.log('--------------------------');
	console.log('\n');

	return joinChannel.processJoinChannel();
}, (err) => {
	console.log('\n');
	console.log('-------------------------');
	console.log('CHANNEL CREATION FAILED:', err);
	console.log('-------------------------');
	console.log('\n');
	process.exit(1);
})
// Join peers to the channel created above
.then(() => {
	console.log('\n');
	console.log('----------------------');
	console.log('CHANNEL JOIN COMPLETE');
	console.log('----------------------');
	console.log('\n');

	return installCC.installChaincode(Constants.CHAINCODE_PATH, Constants.CHAINCODE_VERSION);
}, (err) => {
	console.log('\n');
	console.log('---------------------');
	console.log('CHANNEL JOIN FAILED:', err);
	console.log('---------------------');
	console.log('\n');
	process.exit(1);
})
// Install chaincode on the channel on all peers
.then(() => {
	console.log('\n');
	console.log('---------------------------');
	console.log('CHAINCODE INSTALL COMPLETE');
	console.log('---------------------------');
	console.log('\n');

	return instantiateCC.instantiateOrUpgradeChaincode(
		Constants.APPDEV_ORG,
		Constants.CHAINCODE_PATH,
		Constants.CHAINCODE_VERSION,
		'init',
		["1000", "1111", "1111", "0.1", "1000", "2222", "2222", "1.00", "2020-06-01", "1000"],
		false
	);
}, (err) => {
	console.log('\n');
	console.log('--------------------------');
	console.log('CHAINCODE INSTALL FAILED:', err);
	console.log('--------------------------');
	console.log('\n');
	process.exit(1);
})
// Instantiate chaincode on the channel on all peers
.then(() => {
	console.log('\n');
	console.log('-------------------------------');
	console.log('CHAINCODE INSTANTIATE COMPLETE');
	console.log('-------------------------------');
	console.log('\n');
	ClientUtils.txEventsCleanup();

}, (err) => {
	console.log('\n');
	console.log('------------------------------');
	console.log('CHAINCODE INSTANTIATE FAILED:', err);
	console.log('------------------------------');
	console.log('\n');
	process.exit(1);
});

process.on('uncaughtException', err => {
	console.error(err);
	joinChannel.joinEventsCleanup();
});

process.on('unhandledRejection', err => {
	console.error(err);
	joinChannel.joinEventsCleanup();
});

process.on('exit', () => {
	joinChannel.joinEventsCleanup();
	ClientUtils.txEventsCleanup();
});
