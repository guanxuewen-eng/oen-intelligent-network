package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xinewang/oen/internal/config"
	cors "github.com/gin-contrib/cors"
)

func CORS(cfg *config.CORSConfig) gin.HandlerFunc {
	c := cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Duration(cfg.MaxAge) * time.Second,
	}
	return cors.New(c)
}
