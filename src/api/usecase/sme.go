package usecase

import (
	"credit-line-api/src/api/entity"
	"fmt"
	"strconv"
)

type SME struct{}

func (s *SME) Get(b *entity.Business) float64 {
	fifthOfRevenue := b.MonthlyRevenue / 5
	res, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", fifthOfRevenue), 64)
	return res
}
