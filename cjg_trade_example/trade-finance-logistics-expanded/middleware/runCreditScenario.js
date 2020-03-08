/*

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

var tradeID = '2ks89j9';

/////////////////////////////////
// INVOKE AND QUERY OPERATIONS //
/////////////////////////////////

// INVOKE: acceptTrade (Exporter)
invokeCC.invokeChaincode(Constants.EXPORTER_ORG, Constants.CHAINCODE_UPGRADE_VERSION, 'acceptTrade', [tradeID], 'Exporter')
    .then(() => {
        console.log('\n');
        console.log('------------------------------');
        console.log('CHAINCODE INVOCATION COMPLETE');
        console.log('acceptTrade SUCCEEDED');
        console.log('------------------------------');
        console.log('\n');

        // QUERY: getTradeStatus (Importer)
        return queryCC.queryChaincode(Constants.IMPORTER_ORG, Constants.CHAINCODE_VERSION, 'getTradeStatus', [tradeID], 'Importer');
    }, (err) => {
        console.log('\n');
        console.log('-----------------------------');
        console.log('CHAINCODE INVOCATION FAILED:', err);
        console.log('acceptTrade FAILED');
        console.log('-----------------------------');
        console.log('\n');
        process.exit(1);
    })
    .then((result) => {
        console.log('\n');
        console.log('-------------------------');
        console.log('CHAINCODE QUERY COMPLETE');
        console.log('getTradeStatus VALUE:', result);
        console.log('-------------------------');
        console.log('\n');

        // INVOKE: requestLC (Importer)
        return invokeCC.invokeChaincode(Constants.IMPORTER_ORG, Constants.CHAINCODE_UPGRADE_VERSION, 'requestLC', [tradeID], 'Importer');
    }, (err) => {
        console.log('\n');
        console.log('------------------------');
        console.log('CHAINCODE QUERY FAILED:', err);
        console.log('getTradeStatus FAILED');
        console.log('------------------------');
        console.log('\n');
        process.exit(1);
    })
    .then(() => {
        console.log('\n');
        console.log('------------------------------');
        console.log('CHAINCODE INVOCATION COMPLETE');
        console.log('requestLC SUCCEEDED');
        console.log('------------------------------');
        console.log('\n');

        // QUERY: getLCStatus (Importer)
        return queryCC.queryChaincode(Constants.IMPORTER_ORG, Constants.CHAINCODE_VERSION, 'getLCStatus', [tradeID], 'Importer');
    }, (err) => {
        console.log('\n');
        console.log('-----------------------------');
        console.log('CHAINCODE INVOCATION FAILED:', err);
        console.log('requestLC FAILED');
        console.log('-----------------------------');
        console.log('\n');
        process.exit(1);
    })
    .then((result) => {
        console.log('\n');
        console.log('-------------------------');
        console.log('CHAINCODE QUERY COMPLETE');
        console.log('getLCStatus VALUE:', result);
        console.log('-------------------------');
        console.log('\n');

        // INVOKE: issueLC (Importer's Bank)
        return invokeCC.invokeChaincode(Constants.IMPORTER_ORG, Constants.CHAINCODE_UPGRADE_VERSION, 'issueLC', [tradeID, 'lc8349', '12/31/2018', 'E/L', 'B/L'], 'ImportersBank');
    }, (err) => {
        console.log('\n');
        console.log('------------------------');
        console.log('CHAINCODE QUERY FAILED:', err);
        console.log('getLCStatus FAILED');
        console.log('------------------------');
        console.log('\n');
        process.exit(1);
    })
    .then(() => {
        console.log('\n');
        console.log('------------------------------');
        console.log('CHAINCODE INVOCATION COMPLETE');
        console.log('issueLC SUCCEEDED');
        console.log('------------------------------');
        console.log('\n');

        // QUERY: getLCStatus (Importer's Bank)
        return queryCC.queryChaincode(Constants.IMPORTER_ORG, Constants.CHAINCODE_UPGRADE_VERSION, 'getLCStatus', [tradeID], 'ImportersBank');
    }, (err) => {
        console.log('\n');
        console.log('-----------------------------');
        console.log('CHAINCODE INVOCATION FAILED:', err);
        console.log('issueLC FAILED');
        console.log('-----------------------------');
        console.log('\n');
        process.exit(1);
    })
    .then((result) => {
        console.log('\n');
        console.log('-------------------------');
        console.log('CHAINCODE QUERY COMPLETE');
        console.log('getLCStatus VALUE:', result);
        console.log('-------------------------');
        console.log('\n');

        // INVOKE: acceptLC (Exporter's Bank)
        return invokeCC.invokeChaincode(Constants.EXPORTER_ORG, Constants.CHAINCODE_UPGRADE_VERSION, 'acceptLC', [tradeID], 'ExportersBank');
    }, (err) => {
        console.log('\n');
        console.log('------------------------');
        console.log('CHAINCODE QUERY FAILED:', err);
        console.log('getLCStatus FAILED');
        console.log('------------------------');
        console.log('\n');
        process.exit(1);
    })
    .then(() => {
        console.log('\n');
        console.log('------------------------------');
        console.log('CHAINCODE INVOCATION COMPLETE');
        console.log('acceptLC SUCCEEDED');
        console.log('------------------------------');
        console.log('\n');

        // QUERY: getLCStatus (Exporter's Bank)
        return queryCC.queryChaincode(Constants.EXPORTER_ORG, Constants.CHAINCODE_UPGRADE_VERSION, 'getLCStatus', [tradeID], 'ExportersBank');
    }, (err) => {
        console.log('\n');
        console.log('-----------------------------');
        console.log('CHAINCODE INVOCATION FAILED:', err);
        console.log('acceptLC FAILED');
        console.log('-----------------------------');
        console.log('\n');
        process.exit(1);
    })
    .then((result) => {
        console.log('\n');
        console.log('-------------------------');
        console.log('CHAINCODE QUERY COMPLETE');
        console.log('getLCStatus VALUE:', result);
        console.log('-------------------------');
        console.log('\n');

        /*
        NEW CODE: Scenario of acquiring a credit line
         */

        // INVOKE: acceptLC (Exporter's Bank)
        return invokeCC.invokeChaincode(Constants.EXPORTER_ORG, Constants.CHAINCODE_UPGRADE_VERSION, 'getCreditLine', [tradeID, "Lender"], 'Exporter');
    }, (err) => {
        console.log('\n');
        console.log('------------------------');
        console.log('CHAINCODE QUERY FAILED:', err);
        console.log('getCreditLine FAILED');
        console.log('------------------------');
        console.log('\n');
        process.exit(1);
    })
    .then(() => {
        console.log('\n');
        console.log('------------------------------');
        console.log('CHAINCODE INVOCATION COMPLETE');
        console.log('getCreditLine SUCCEEDED');
        console.log('------------------------------');
        console.log('\n');

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
