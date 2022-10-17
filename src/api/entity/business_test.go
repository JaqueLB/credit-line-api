package entity

import (
	"reflect"
	"testing"
	"time"
)

func TestNewBusiness(t *testing.T) {
	testDate, _ := time.Parse("2006-01-02T15:04:05.00Z07:00", "2021-07-19T16:32:59.860Z")
	type args struct {
		params *CreditLineInput
	}
	tests := []struct {
		name string
		args args
		want *Business
	}{
		{
			"Should convert the input into Business",
			args{
				params: &CreditLineInput{
					FoundingType:      "SME",
					CashBalance:       100.10,
					MonthlyRevenue:    1000.89,
					RequestedValue:    50.15,
					RequestedDateTime: &testDate,
				},
			},
			&Business{
				Type:           FoundingTypeSME,
				CashBalance:    100.10,
				MonthlyRevenue: 1000.89,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBusiness(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBusiness() = %v, want %v", got, tt.want)
			}
		})
	}
}
