package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CheckHandler struct {
	router *gin.Engine
}

func NewCheckHandler(router *gin.Engine) *CheckHandler {
	return &CheckHandler{
		router: router,
	}
}

func (h *CheckHandler) Register() {

	h.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, &ServiceInfoResponse{Service: "qyrodns"})
	})

	h.router.GET("/health", func(c *gin.Context) {
		c.JSONP(http.StatusOK, &CheckResponse{Status: "ok"})
	})
}
