package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SupabaseAuthMiddleware validates the JWT token issued by Supabase
func SupabaseAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		jwtSecret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set claims in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["sub"])
			c.Set("email", claims["email"])
			c.Set("app_metadata", claims["app_metadata"]) // To check roles later if using custom claims
		}

		c.Next()
	}
}

// RequireRoleMiddleware checks if the user has a specific role (requires SupabaseAuthMiddleware to run first)
func RequireRoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		metadata, exists := c.Get("app_metadata")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: No metadata found"})
			return
		}

		metaMap, ok := metadata.(map[string]interface{})
		if !ok || metaMap["role"] != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient role"})
			return
		}

		c.Next()
	}
}
