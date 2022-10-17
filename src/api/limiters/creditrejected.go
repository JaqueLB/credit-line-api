package limiters

import (
	"time"

	"golang.org/x/time/rate"
)

type CLRejectedLimiter struct{}

func (l *CLRejectedLimiter) Get() *rate.Limiter {
	return rate.NewLimiter(rate.Every(2*time.Minute), 1) // 1 req/30 s
}
