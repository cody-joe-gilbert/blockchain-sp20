'use strict';

var Constants = require('./constants.js');
var ClientUtils = require('./clientUtils.js');
var joinChannel = require('./join-channel.js');
var invokeCC = require('./invoke-chaincode.js');
var queryCC = require('./query-chaincode.js');

var inputArgs = process.argv;
var orgName = inputArgs[2]
var funcName = inputArgs[3]
var argList = inputArgs[4]
var userName = inputArgs[5]
queryCC.queryChaincode(
    orgName,
    Constants.CHAINCODE_VERSION,
    funcName,
    argList,
    userName)
    .then((result) => {
        console.log(result);
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