package okserver

import (
	"github.com/gin-gonic/gin"
	"github.com/okgotool/okgoweb/okmodel/okresponse"
)

var (
	Healthcheck = &HealthcheckType{}
)

type (
	HealthcheckType struct {
	}
)

func (w *HealthcheckType) AddApis(router *gin.Engine) {
	router.GET("/hello", w.HelloHandle)
	router.GET("/healthz", w.HealthzHandle)
}

// HelloHandle ：
// Swagger doc refer: https://github.com/swaggo/swag
// @Summary Hello API
// @Description Hello test API
// @Success 200 {object} okresponse.Success "{"msg":"ok"}"
// @Router /hello [get]
func (w *HealthcheckType) HelloHandle(c *gin.Context) {

	c.JSON(okresponse.OKCode, &okresponse.Success{
		Code: okresponse.StatusOK,
		Msg:  "ok",
	})
}

// HealthzHandle ：
// Swagger doc refer: https://github.com/swaggo/swag
// @Summary Healthz API
// @Description Health check API
// @Success 200 {object} okresponse.Success "{\"msg\":\"ok\"}"
// @Router /healthz [get]
func (w *HealthcheckType) HealthzHandle(c *gin.Context) {
	c.JSON(okresponse.OKCode, &okresponse.Success{
		Code: okresponse.StatusOK,
		Msg:  "ok",
	})
}
