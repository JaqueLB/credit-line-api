package entity

import (
	"golang.org/x/time/rate"
)

type ILimiter interface {
	Get() *rate.Limiter
}
