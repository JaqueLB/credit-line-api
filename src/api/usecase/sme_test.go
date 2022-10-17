package usecase

import (
	"credit-line-api/src/api/entity"
	"testing"
)

func TestSME_Get(t *testing.T) {
	type args struct {
		b *entity.Business
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"Should return the right recommended value for SMEs",
			args{
				&entity.Business{
					Type:           entity.FoundingTypeSME,
					CashBalance:    900.99,
					MonthlyRevenue: 5000.20,
				},
			},
			1000.04,
		},
		{
			"Should return the recommended value for SMEs based on Monthly Revenue, even if this is lower",
			args{
				&entity.Business{
					Type:           entity.FoundingTypeSME,
					CashBalance:    900.99,
					MonthlyRevenue: 500.20,
				},
			},
			100.04,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SME{}
			if got := s.Get(tt.args.b); got != tt.want {
				t.Errorf("SME.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
