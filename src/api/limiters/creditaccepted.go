package limiters

import (
	"time"

	"golang.org/x/time/rate"
)

type CLAcceptedLimiter struct{}

func (l *CLAcceptedLimiter) Get() *rate.Limiter {
	return rate.NewLimiter(rate.Every(time.Minute), 2) // 2 reqs/2 min
}
