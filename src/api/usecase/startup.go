package usecase

import "credit-line-api/src/api/entity"

type Startup struct{}

func (s *Startup) Get(b *entity.Business) float64 {
	thirdOfCash := b.CashBalance / 3
	fifthOfRevenue := b.MonthlyRevenue / 5

	if thirdOfCash > fifthOfRevenue {
		return thirdOfCash
	}

	return fifthOfRevenue
}
