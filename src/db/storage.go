package db

import (
	"credit-line-api/src/api/entity"
)

type IStorage interface {
	Get(key int) *entity.CreditLineResponse
	Set(key int, value *entity.CreditLineResponse) bool
}
