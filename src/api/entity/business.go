package entity

type FoundingType string

const (
	FoundingTypeSME     FoundingType = "sme"
	FoundingTypeStartup FoundingType = "startup"
)

type Business struct {
	Type           FoundingType
	CashBalance    float64
	MonthlyRevenue float64
}

func NewBusiness(params *CreditLineInput) *Business {
	return &Business{
		Type:           FoundingType(params.FoundingType),
		CashBalance:    params.CashBalance,
		MonthlyRevenue: params.MonthlyRevenue,
	}
}
