package usecase

import (
	"credit-line-api/src/api/entity"
	"fmt"
	"strconv"
)

type Startup struct{}

func (s *Startup) Get(b *entity.Business) float64 {
	thirdOfCash := b.CashBalance / 3
	fifthOfRevenue := b.MonthlyRevenue / 5
	var res float64

	if thirdOfCash > fifthOfRevenue {
		res, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", thirdOfCash), 64)
	} else {
		res, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", fifthOfRevenue), 64)
	}

	return res
}
