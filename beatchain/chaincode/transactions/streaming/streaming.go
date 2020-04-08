package streaming

var StreamingVariable = "streaming"

type Contract struct {
	CreatorId			int			`json:"creatorId"`
	AppDevId			int			`json:"appDevId"`
	ProductId			int			`json:"productId"`
	DevShareFraction	float64		`json:"devShareFraction"`
}

type Product struct {
	ProductId			int			`json:"productId"`
	CreatorId			int			`json:"creatorId"`
	ProductName			string		`json:"productName"`
	TotalListens		int			`json:"totalListens"`
	TotalMetrics1		string		`json:"totalMetrics1"`
	Totalmetrics2		string		`json:"totalMetrics2"`
	UnpaidListens		int			`json:"unpaidListens"`
	UnpaidMetrics		int			`json:"unpaidMetrics"`
}