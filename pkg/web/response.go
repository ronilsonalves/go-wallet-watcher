package web

import (
	"github.com/gin-gonic/gin"
	"time"
)

type errResponse struct {
	StatusCode int    `json:"status-code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	TimeStamp  string `json:"time-stamp"`
}

type PageableResponse struct {
	Page  string      `json:"page"`
	Items string      `json:"items"`
	Data  interface{} `json:"data"`
}

func BadResponse(ctx *gin.Context, stsCode int, status, message string) {
	ctx.AbortWithStatusJSON(stsCode, errResponse{
		StatusCode: stsCode,
		Status:     status,
		Message:    message,
		TimeStamp:  time.Now().String(),
	})
}

func OKResponse(ctx *gin.Context, stsCode int, data interface{}) {
	ctx.JSON(stsCode, data)
}
