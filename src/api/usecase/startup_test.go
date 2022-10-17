package usecase

import (
	"credit-line-api/src/api/entity"
	"testing"
)

func TestStartup_Get(t *testing.T) {
	type args struct {
		b *entity.Business
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"Should return the recommended value for Startups based on the Cash Balance",
			args{
				&entity.Business{
					Type:           entity.FoundingTypeStartup,
					CashBalance:    900.99,
					MonthlyRevenue: 500.20,
				},
			},
			300.33,
		},
		{
			"Should return the recommended value for Startups based on the Monthly Revenue",
			args{
				&entity.Business{
					Type:           entity.FoundingTypeStartup,
					CashBalance:    900.99,
					MonthlyRevenue: 5000.20,
				},
			},
			1000.04,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Startup{}
			if got := s.Get(tt.args.b); got != tt.want {
				t.Errorf("Startup.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
