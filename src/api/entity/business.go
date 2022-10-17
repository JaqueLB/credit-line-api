package entity

import "strings"

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
	foundingType := strings.ToLower(params.FoundingType)

	return &Business{
		Type:           FoundingType(foundingType),
		CashBalance:    params.CashBalance,
		MonthlyRevenue: params.MonthlyRevenue,
	}
}
