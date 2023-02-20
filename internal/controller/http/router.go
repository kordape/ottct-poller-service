// Package http implements routing paths. Each services in own file.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	v1 "github.com/kordape/ottct-main-service/internal/controller/http/v1"
	"github.com/kordape/ottct-main-service/pkg/logger"
)

func NewRouter(handler *gin.Engine, l logger.Interface) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	h := handler.Group("/v1")
	{
		v1.NewRoutes(h, l)
	}
}
