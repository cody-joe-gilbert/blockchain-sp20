package admin

var AdminVariable = "admin"

type CreatorRecord struct {
	CreatorId		int			`json:"creatorId"`
	BankId			int			`json:"bankId"`
}

type AppDevRecord struct {
	AppDevId		int			`json:"appDevId"`
	BankId			int			`json:"bankId"`
	AdminFeeFrac	float64		`json:"adminFeeFrac"`
}

type CustomerRecord struct {
	CustomerId				int			`json:"costumerId"`
	AppDevId				int			`json:"appDevId"`
	BankId					int			`json:"bankId"`
	SubscriptionFee			float32		`json:"subscriptionFee"`
	SubscriptionDueDate		string		`json:"subscriptionDueDate"`
	QueuedSong				int			`json:"productId"`
	PreviousSong			int			`json:"productId"`
}