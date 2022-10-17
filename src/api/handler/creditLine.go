package handler

import (
	"credit-line-api/src/api/entity"
	"credit-line-api/src/api/limiters"
	"credit-line-api/src/api/usecase"
	"credit-line-api/src/db"
	"net/http"
	"strconv"
	"strings"

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

	foundingType := strings.ToLower(body.FoundingType)
	if foundingType != string(entity.FoundingTypeSME) &&
		foundingType != string(entity.FoundingTypeStartup) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unsupported founding type",
		})
		return
	}

	if body.CashBalance <= 0 ||
		body.MonthlyRevenue <= 0 ||
		body.RequestedValue <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid credit line request",
		})
		return
	}

	// check if credit was approved before
	res := h.Storage.Get(clientID)
	if res != nil {
		if res.Accepted == true {
			limiter := &limiters.CLAcceptedLimiter{}
			if !limiter.Get().Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"message": http.StatusText(http.StatusTooManyRequests),
				})
				return
			}

			c.JSON(http.StatusOK, res)
			return
		} else {
			if res.RejectedCount >= 3 {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "A sales agent will contact you",
				})
				return
			}
		}
	} else {
		res = &entity.CreditLineResponse{
			Accepted:      false,
			RejectedCount: 0,
		}
	}

	// calculate the recommended credit line
	business := entity.NewBusiness(body)
	recommendedValue := getRecommendedCreditLine(business)

	// check if credit is approved
	if recommendedValue > body.RequestedValue {
		res.Accepted = true
		res.ApprovedValue = body.RequestedValue
		h.Storage.Set(clientID, res)
	} else {
		// in case credit is rejected, limit the access
		res.RejectedCount = res.RejectedCount + 1
		h.Storage.Set(clientID, res)
		limiter := &limiters.CLRejectedLimiter{}
		if !limiter.Get().Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": http.StatusText(http.StatusTooManyRequests),
			})
			return
		}
	}

	c.JSON(http.StatusOK, res)
	return
}
