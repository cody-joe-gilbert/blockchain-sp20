package utils

import (
	"time"
)

const UNIQUE_ID_KEY = "UNIQUE_ID_VARIABLE"
const UNIQUE_STARTING_ID = "100000000"

const BEATCHAIN_ADMIN_MSP = "BeatchainMSP"
const BEATCHAIN_ADMIN_CA = "ca.admin.beatchain.com"
const BEATCHAIN_ADMIN_BANK_ACCOUNT_ID = "1"
// to be done:  need to factor in changes for bank account inUse bool - Arun
// Authorization constants
const CUSTOMER_MSP = "CustomerMSP"
const CUSTOMER_CA = "ca.customerorg.beatchain.com"
const APPDEV_MSP = "AppDevMSP"
const APPDEV_CA = "ca.appdevorg.beatchain.com"
const CREATOR_MSP = "CreatorMSP"
const CREATOR_CA = "ca.creatororg.beatchain.com"

// Key Constants
const KEY_OBJECT_FORMAT = "object~id"
const CUSTOMER_RECORD_KEY_PREFIX = "CustomerRecord"
const CREATOR_RECORD_KEY_PREFIX = "CreatorRecord"
const CONTRACT_KEY_PREFIX = "Contract"
const BANK_ACCOUNT_KEY_PREFIX = "BankAccount"
const APPDEV_RECORD_KEY_PREFIX = "AppDevRecord"
const PRODUCT_KEY_PREFIX = "Product"

// Test constants
const BEATCHAIN_ADMIN_BALANCE = "1000"

const TEST_APPDEV_ID = "1111"
const TEST_APPDEV_BA_ID = "1111"
const TEST_APPDEV_DEVSHARE = "0.1"
const TEST_APPDEV_BA_BALANCE = "1000"

const TEST_CUSTOMER_ID = "2222"
const TEST_CUSTOMER_BA_ID = "2222"
const TEST_CUSTOMER_SUBFEE = "1.00"
const TEST_CUSTOMER_SUB_DUE_DATE = "2020-06-01"
const TEST_CUSTOMER_BA_BALANCE = "1000"

const TEST_CREATOR_ID = "3333"
const TEST_CREATOR_BA_ID = "3333"
const TEST_CREATOR_BA_BALANCE = "1000"

const TEST_PRODUCT_ID = "4444"
const TEST_PRODUCT_NAME = "Test Product"
const TEST_PRODUCT_TOTAL_LISTENS = "5"
const TEST_PRODUCT_UNREN_LISTENS = "3"
const TEST_PRODUCT_TOTAL_METRICS = "7"
const TEST_PRODUCT_UNREN_METRICS = "4"
const TEST_PRODUCT_ADD_METRICS = "0"
const TEST_PRODUCT_ACTIVE = "true"

const TEST_CONTRACT_PPS = "0.01"
const TEST_CONTRACT_STATUS = "true"

type Transaction struct {
	/*
		Defines the standard format of the user's transaction info. Used to prevent having to
		replicate boilerplate code on each transaction function
	*/
	CalledFunction    string
	CreatorId         string
	CreatorOrg        string
	CreatorCertIssuer string
	CreatorAdmin      bool
	Args              []string
	TestMode		  bool
}

type CustomerRecord struct {
	/*
		Defines a single customer for a single AppDev on the ledger
	*/
	Id                  string    `json:"id"`
	AppDevId            string    `json:"appdevid"`
	BankAccountId       string    `json:"bankaccountid"`
	SubscriptionFee     float32   `json:"subscriptionfee"`
	SubscriptionDueDate time.Time `json:"subscriptionduedate"`
	QueuedSong          string    `json:"queuedsong"`
	PreviousSong        string    `json:"previoussong"`
}

type CreatorRecord struct {
	/*
		Defines a single Creator
	*/
	Id                  string    `json:"id"`
	BankAccountId       string    `json:"bankaccountid"`
}

type BankAccount struct {
	/*
		Defines a bank account on the ledger
	*/
	Id      string  `json:"id"`
	Balance float32 `json:"balance"`
	InUse	bool 	`json:"inUse"` //if true cant be assigned to a new entity. can be assigned to only one person
}

type AppDevRecord struct {
	/*
		Defines a single AppDev record on the ledger
	*/
	Id            string  `json:"id"`
	BankAccountId string  `json:"bankaccountid"`
	AdminFeeFrac  float32 `json:"adminfeefrac"`
}

type Contract struct {
	/*
		Defines a Contract record on the ledger
	*/
	CreatorId           string  `json:"creatorid"`
	AppDevId            string  `json:"appdevid"`
	ProductId           string  `json:"productid"`
	CreatorPayPerStream float32 `json:"creatorpayperstream"`
	Status              string  `json:"contractstatus"`
}

type Product struct {
	/*
		Defines a product and its attributes on the ledger
	*/

	Id                   string `json:"id"`
	CreatorId            string `json:"id"`
	ProductName          string `json:"productName"`
	TotalListens         int64  `json:"totalListens"`
	UnRenumeratedListens int64  `json:"unRenumeratedListens"`
	TotalMetrics         int64  `json:"totalMetrics"`
	UnRenumeratedMetrics int64  `json:"unRenumeratedMetrics"`
	AdditionalMetrics    int64  `json:"additionalMetrics"`
	IsActive             bool   `json:"isActive"`
}
