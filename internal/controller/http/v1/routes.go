package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/kordape/ottct-main-service/pkg/logger"
)

type routes struct {
	l logger.Interface
}

func NewRoutes(handler *gin.RouterGroup, l logger.Interface) {
	r := &routes{l}

	h := handler.Group("/echo")
	{
		h.GET("/", r.echoHandler)
	}
}
