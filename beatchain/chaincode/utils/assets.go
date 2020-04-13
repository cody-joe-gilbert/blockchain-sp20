package utils

import (
	"time"
)

const BEATCHAIN_ADMIN_BANK_ACCOUNT_ID = "0000001"

type Transaction struct {
	/*
		Defines the standard format of the user's transaction info. Used to prevent having to
		replicate boilerplate code on each transaction function
	*/
	CalledFunction    string
	CreatorId         string
	CreatorOrg        string
	CreatorCertIssuer string
	Args              []string
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

type BankAccount struct {
	/*
		Defines a bank account on the ledger
	*/
	Id      string  `json:"id"`
	Balance float32 `json:"balance"`
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
	CreatorId           string  `json:"id"`
	AppDevId            string  `json:"id"`
	ProductId           string  `json:"id"`
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
