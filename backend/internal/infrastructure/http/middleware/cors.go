package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware permite que el frontend (Vercel) llame al backend desde otro
// origen. Sin esto el navegador bloquea las peticiones cross-origin y las
// respuestas llegan como "NetworkError". Usamos Bearer token (no cookies),
// así que no necesitamos Allow-Credentials.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if isAllowedOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	// Producción en Vercel (incluye los preview deploys *.vercel.app) y dev local.
	// ponytail: lista fija; si el dominio cambia, se agrega aquí.
	if strings.HasSuffix(origin, ".vercel.app") {
		return true
	}
	switch origin {
	case "http://localhost:4200", "http://localhost:3000":
		return true
	}
	return false
}
