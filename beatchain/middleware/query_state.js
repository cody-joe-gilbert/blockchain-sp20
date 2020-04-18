'use strict';

var Constants = require('./constants.js');
var ClientUtils = require('./clientUtils.js');
var joinChannel = require('./join-channel.js');
var invokeCC = require('./invoke-chaincode.js');
var queryCC = require('./query-chaincode.js');


queryCC.queryChaincode(
    Constants.APPDEV_ORG,
    Constants.CHAINCODE_VERSION,
    'ListBankAccounts',
    [],
    'admin')
    .then((result) => {
        console.log('\n');
        console.log('-------------------------');
        console.log('CHAINCODE QUERY COMPLETE');
        console.log('VALUE:', result);
        console.log('-------------------------');
        console.log('\n');
        ClientUtils.txEventsCleanup();
    }, (err) => {
        console.log('\n');
        console.log('------------------------');
        console.log('CHAINCODE QUERY FAILED:', err);
        console.log('getTradeStatus FAILED');
        console.log('------------------------');
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
});
