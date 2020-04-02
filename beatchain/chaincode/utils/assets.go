package utils

import "time"

const BEATCHAIN_ADMIN_BANK_ACCOUNT_ID = "0000001"

type Transaction struct {
	/*
	Defines the standard format of the user's transaction info. Used to prevent having to
	replicate boilerplate code on each transaction function
	 */
	CalledFunction string
	CreatorId string
	CreatorOrg string
	CreatorCertIssuer string
	Args []string
}

type CustomerRecord struct {
	/*
	Defines a single customer for a single AppDev on the ledger
	 */
	Id							string		`json:"id"`
	AppDevId					string		`json:"appdevid"`
	BankAccountId				string		`json:"bankaccountid"`
	SubscriptionFee				float32		`json:"subscriptionfee"`
	SubscriptionDueDate			time.Time	`json:"subscriptionduedate"`
	QueuedSong					string		`json:"qeuedsong"`
	PreviousSong				string		`json:"previoussong"`
}

type BankAccount struct {
	/*
	Defines a bank account on the ledger
	 */
	Id							string		`json:"id"`
	Balance						float32		`json:"balance"`
}

type AppDevRecord struct {
	/*
	Defines a single AppDev record on the ledger
	*/
	Id							string		`json:"id"`
	BankAccountId				string		`json:"bankaccountid"`
	AdminFeeFrac				float32		`json:"adminfeefrac"`
}
