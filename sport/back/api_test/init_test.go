package api_test

import (
	"io/ioutil"
	"sport/api"
	"sport/config"

	"github.com/gin-gonic/gin"
)

var apiCtx *api.Ctx

func init() {
	apiCtx = api.NewApi(config.NewLocalCtx())

	apiCtx.JwtGroup.GET("/stat", func(c *gin.Context) {
		c.AbortWithStatus(204)
	})

	apiCtx.AnonGroup.POST("/read_body", func(c *gin.Context) {
		if _, err := ioutil.ReadAll(c.Request.Body); err != nil {
			if !c.IsAborted() {
				c.AbortWithError(400, err)
			}
		} else {
			c.AbortWithStatus(204)
		}
	})
}
