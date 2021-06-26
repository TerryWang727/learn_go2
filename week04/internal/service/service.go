package service

import (
	"Week04/internal/biz"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	apiv1 := r.Group("/api/v1")

	{
		apiv1.GET("/tags", biz.GetTags)
	}

	return r
}