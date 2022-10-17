package db

import (
	"credit-line-api/src/api/entity"
)

type LocalStorage struct {
	Items map[int]*entity.CreditLineResponse
}

func (s *LocalStorage) Get(k int) *entity.CreditLineResponse {
	return s.Items[k]
}

func (s *LocalStorage) Set(k int, v *entity.CreditLineResponse) bool {
	s.Items[k] = v
	return true
}
