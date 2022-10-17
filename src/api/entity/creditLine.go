package entity

import "time"

type CreditLineInput struct {
	FoundingType      string     `json:"founding_type"`
	CashBalance       float64    `json:"cash_balance"`
	MonthlyRevenue    float64    `json:"monthly_revenue"`
	RequestedValue    float64    `json:"requested_value"`
	RequestedDateTime *time.Time `json:"requested_datetime"`
}

type CreditLineResponse struct {
	Accepted      bool    `json:"accepted"`
	ApprovedValue float64 `json:"approved_value"`
	RejectedCount int     `json:"-"`
}

type ICreditLineCalculator interface {
	Get(*Business) float64
}
