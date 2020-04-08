package banking

var BankingVariable = "banking"

type BankAccount struct {
	BankId			int			`json:"bankId"`
	Balance			int			`json:"balance"`
}