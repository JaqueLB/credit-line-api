package handler

import (
	"credit-line-api/src/api/entity"
	"credit-line-api/src/api/limiters"
	"credit-line-api/src/api/usecase"
	"credit-line-api/src/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreditLineHandler struct {
	Storage db.IStorage
}

func getCalculator(t entity.FoundingType) entity.ICreditLineCalculator {
	if t == entity.FoundingTypeSME {
		return &usecase.SME{}
	}
	return &usecase.Startup{}
}

func getRecommendedCreditLine(b *entity.Business) float64 {
	calculator := getCalculator(b.Type)
	return calculator.Get(b)
}

func (h *CreditLineHandler) Check(c *gin.Context) {
	// validate data
	clientID, err := strconv.Atoi(c.Params.ByName("clientID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid clientID",
		})
		return
	}

	var body *entity.CreditLineInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
		})
		return
	}

	if body.CashBalance == 0 && body.MonthlyRevenue == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid credit line request",
		})
		return
	}

	// check if credit was approved before
	res := h.Storage.Get(clientID)
	if res != nil && res.Accepted {
		limiter := &limiters.CLAcceptedLimiter{}
		if !limiter.Get().Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": http.StatusText(http.StatusTooManyRequests),
			})
			return
		}

		c.JSON(http.StatusOK, res)
		return
	}

	// calculate the recommended credit line
	business := entity.NewBusiness(body)
	recommendedValue := getRecommendedCreditLine(business)

	// check if credit is approved
	if recommendedValue > body.RequestedValue {
		res.Accepted = true
		res.ApprovedValue = body.RequestedValue
		h.Storage.Set(clientID, res)
		c.JSON(http.StatusOK, res)
		return
	}

	// in case credit is rejected, limit the access
	res.Accepted = false
	limiter := &limiters.CLRejectedLimiter{}
	if !limiter.Get().Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"message": http.StatusText(http.StatusTooManyRequests),
		})
		return
	}
	// TODO: check if client was rejected before, after 3 times give an error
}
