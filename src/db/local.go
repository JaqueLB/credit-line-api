package db

import (
	"credit-line-api/src/api/entity"
)

var localItems = make(map[int]*entity.CreditLineResponse)

type LocalStorage struct{}

func (s *LocalStorage) Get(k int) *entity.CreditLineResponse {
	return localItems[k]
}

func (s *LocalStorage) Set(k int, v *entity.CreditLineResponse) bool {
	localItems[k] = v
	return true
}
