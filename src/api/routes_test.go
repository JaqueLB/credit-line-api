package api

import (
	"bytes"
	"credit-line-api/src/api/entity"
	"credit-line-api/src/api/handler"
	"credit-line-api/src/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreditLineHandler_Check_Validations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	request := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(request)
	handlers := &Handlers{
		CreditLineHandler: &handler.CreditLineHandler{
			Storage: &db.LocalStorage{},
		},
	}
	router := CreateRouter(handlers, engine)

	type args struct {
		clientID string
		reqBody  *entity.CreditLineInput
	}
	type want struct {
		res  *entity.CreditLineResponse
		code int
		msg  string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"Should return error when client is not an int",
			args{
				clientID: "1a",
				reqBody:  &entity.CreditLineInput{},
			},
			want{
				res:  nil,
				code: http.StatusBadRequest,
				msg:  "Invalid clientID",
			},
		},
		{
			"Should return error when founding type is empty",
			args{
				clientID: "1",
				reqBody:  &entity.CreditLineInput{},
			},
			want{
				res:  nil,
				code: http.StatusBadRequest,
				msg:  "Unsupported founding type",
			},
		},
		{
			"Should return error when founding type is wrong",
			args{
				clientID: "1",
				reqBody: &entity.CreditLineInput{
					FoundingType: "test",
				},
			},
			want{
				res:  nil,
				code: http.StatusBadRequest,
				msg:  "Unsupported founding type",
			},
		},
		{
			"Should return error when required properties are empty",
			args{
				clientID: "1",
				reqBody: &entity.CreditLineInput{
					FoundingType: "SME",
				},
			},
			want{
				res:  nil,
				code: http.StatusBadRequest,
				msg:  "Invalid credit line request",
			},
		},
		{
			"Should return error when cash is zero",
			args{
				clientID: "1",
				reqBody: &entity.CreditLineInput{
					FoundingType:   "SME",
					CashBalance:    0,
					RequestedValue: 2,
					MonthlyRevenue: 3,
				},
			},
			want{
				res:  nil,
				code: http.StatusBadRequest,
				msg:  "Invalid credit line request",
			},
		},
		{
			"Should return error when cash is negative",
			args{
				clientID: "1",
				reqBody: &entity.CreditLineInput{
					FoundingType:   "SME",
					CashBalance:    -1,
					MonthlyRevenue: 2,
					RequestedValue: 3,
				},
			},
			want{
				res:  nil,
				code: http.StatusBadRequest,
				msg:  "Invalid credit line request",
			},
		},
		{
			"Should return error when revenue is zero",
			args{
				clientID: "1",
				reqBody: &entity.CreditLineInput{
					FoundingType:   "SME",
					CashBalance:    1,
					MonthlyRevenue: 0,
					RequestedValue: 3,
				},
			},
			want{
				res:  nil,
				code: http.StatusBadRequest,
				msg:  "Invalid credit line request",
			},
		},
		{
			"Should return error when revenue is negative",
			args{
				clientID: "1",
				reqBody: &entity.CreditLineInput{
					FoundingType:   "SME",
					CashBalance:    1,
					MonthlyRevenue: -2,
					RequestedValue: 3,
				},
			},
			want{
				res:  nil,
				code: http.StatusBadRequest,
				msg:  "Invalid credit line request",
			},
		},
		{
			"Should return error when requested value is zero",
			args{
				clientID: "1",
				reqBody: &entity.CreditLineInput{
					FoundingType:   "SME",
					CashBalance:    1,
					MonthlyRevenue: 2,
					RequestedValue: 0,
				},
			},
			want{
				res:  nil,
				code: http.StatusBadRequest,
				msg:  "Invalid credit line request",
			},
		},
		{
			"Should return error when requested value is negative",
			args{
				clientID: "1",
				reqBody: &entity.CreditLineInput{
					FoundingType:   "SME",
					CashBalance:    1,
					MonthlyRevenue: 2,
					RequestedValue: -3,
				},
			},
			want{
				res:  nil,
				code: http.StatusBadRequest,
				msg:  "Invalid credit line request",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			byts, _ := json.Marshal(tt.args.reqBody)
			req, _ := http.NewRequest("POST", fmt.Sprintf("/credit-line/%s", tt.args.clientID), bytes.NewBuffer(byts))
			req.Header.Add("Content-Type", "application/json;charset=utf-8")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			jsonResp, _ := ioutil.ReadAll(resp.Body)

			assert.Equal(t, tt.want.code, resp.Code)
			assert.Contains(t, string(jsonResp), tt.want.msg)
		})
	}
}

func TestCreditLineHandler_Check_RecurringClient(t *testing.T) {
	type args struct {
		clientID string
		reqBody  *entity.CreditLineInput
		storage  *db.LocalStorage
	}
	type want struct {
		res  *entity.CreditLineResponse
		code int
		msg  string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"Should return error when the credit was rejected 3 or more times",
			args{
				clientID: "1",
				reqBody: &entity.CreditLineInput{
					FoundingType:   "SME",
					CashBalance:    1,
					MonthlyRevenue: 2,
					RequestedValue: 3,
				},
				storage: &db.LocalStorage{
					Items: map[int]*entity.CreditLineResponse{
						1: {
							Accepted:      false,
							ApprovedValue: 0,
							RejectedCount: 3,
						},
					},
				},
			},
			want{
				res:  nil,
				code: http.StatusBadRequest,
				msg:  "A sales agent will contact you",
			},
		},
		{
			"Should return the same credit line when the credit was approved, even with different values",
			args{
				clientID: "1",
				reqBody: &entity.CreditLineInput{
					FoundingType:   "SME",
					CashBalance:    3000,
					MonthlyRevenue: 5000,
					RequestedValue: 1000,
				},
				storage: &db.LocalStorage{
					Items: map[int]*entity.CreditLineResponse{
						1: {
							Accepted:      true,
							ApprovedValue: 900,
							RejectedCount: 0,
						},
					},
				},
			},
			want{
				res: &entity.CreditLineResponse{
					Accepted:      true,
					ApprovedValue: 900,
				},
				code: http.StatusOK,
				msg:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			request := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(request)
			handlers := &Handlers{
				CreditLineHandler: &handler.CreditLineHandler{
					Storage: tt.args.storage,
				},
			}
			router := CreateRouter(handlers, engine)

			byts, _ := json.Marshal(tt.args.reqBody)
			req, _ := http.NewRequest("POST", fmt.Sprintf("/credit-line/%s", tt.args.clientID), bytes.NewBuffer(byts))
			req.Header.Add("Content-Type", "application/json;charset=utf-8")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			jsonResp, _ := ioutil.ReadAll(resp.Body)
			var converted *entity.CreditLineResponse
			json.Unmarshal(jsonResp, &converted)

			assert.Equal(t, tt.want.code, resp.Code)
			assert.Contains(t, string(jsonResp), tt.want.msg)
			if tt.want.res != nil {
				assert.Equal(t, tt.want.res.Accepted, converted.Accepted)
				assert.Equal(t, tt.want.res.ApprovedValue, converted.ApprovedValue)
			}
		})
	}
}

func TestCreditLineHandler_Check_Calculator(t *testing.T) {
	gin.SetMode(gin.TestMode)
	request := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(request)
	handlers := &Handlers{
		CreditLineHandler: &handler.CreditLineHandler{
			Storage: &db.LocalStorage{
				Items: map[int]*entity.CreditLineResponse{},
			},
		},
	}
	router := CreateRouter(handlers, engine)

	type args struct {
		clientID string
		reqBody  *entity.CreditLineInput
	}
	type want struct {
		res  *entity.CreditLineResponse
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"Should return success even though the credit was rejected",
			args{
				clientID: "1",
				reqBody: &entity.CreditLineInput{
					FoundingType:   "SME",
					CashBalance:    100,
					MonthlyRevenue: 200,
					RequestedValue: 3000,
				},
			},
			want{
				res: &entity.CreditLineResponse{
					Accepted:      false,
					ApprovedValue: 0,
				},
				code: http.StatusOK,
			},
		},
		{
			"Should return success when the credit is approved for SME",
			args{
				clientID: "2",
				reqBody: &entity.CreditLineInput{
					FoundingType:   "smE",
					CashBalance:    3000,
					MonthlyRevenue: 5000,
					RequestedValue: 950,
				},
			},
			want{
				res: &entity.CreditLineResponse{
					Accepted:      true,
					ApprovedValue: 950,
				},
				code: http.StatusOK,
			},
		},
		{
			"Should return success when the credit is approved for Startup",
			args{
				clientID: "3",
				reqBody: &entity.CreditLineInput{
					FoundingType:   "startup",
					CashBalance:    3000,
					MonthlyRevenue: 100,
					RequestedValue: 940,
				},
			},
			want{
				res: &entity.CreditLineResponse{
					Accepted:      true,
					ApprovedValue: 940,
				},
				code: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			byts, _ := json.Marshal(tt.args.reqBody)
			req, _ := http.NewRequest("POST", fmt.Sprintf("/credit-line/%s", tt.args.clientID), bytes.NewBuffer(byts))
			req.Header.Add("Content-Type", "application/json;charset=utf-8")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			jsonResp, _ := ioutil.ReadAll(resp.Body)
			var converted *entity.CreditLineResponse
			json.Unmarshal(jsonResp, &converted)

			assert.Equal(t, tt.want.code, resp.Code)
			assert.Equal(t, tt.want.res.Accepted, converted.Accepted)
			assert.Equal(t, tt.want.res.ApprovedValue, converted.ApprovedValue)
		})
	}
}
