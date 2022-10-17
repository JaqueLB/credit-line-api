package usecase

import "credit-line-api/src/api/entity"

type SME struct{}

func (s *SME) Get(b *entity.Business) float64 {
	return b.MonthlyRevenue / 5
}
