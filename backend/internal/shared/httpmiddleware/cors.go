package httpmiddleware

import (
	"net/http"
	"strings"

	"cartera-app/backend/internal/config"

	"github.com/gin-gonic/gin"
)

func CORS(cfg *config.Config) gin.HandlerFunc {
	allowedOrigins := append([]string(nil), cfg.CORSAllowedOrigins...)

	return func(c *gin.Context) {
		if len(allowedOrigins) == 0 {
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}

			c.Next()
			return
		}

		origin := strings.TrimSpace(c.GetHeader("Origin"))
		allowed := origin != "" && isOriginAllowed(origin, allowedOrigins)

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == http.MethodOptions {
			if origin == "" || allowed {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}

			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		pattern := strings.TrimSpace(allowed)
		if pattern == "" {
			continue
		}

		if pattern == "*" || strings.EqualFold(pattern, origin) {
			return true
		}

		if strings.Contains(pattern, "*") && wildcardMatch(origin, pattern) {
			return true
		}
	}

	return false
}

func wildcardMatch(value, pattern string) bool {
	if pattern == "*" {
		return true
	}

	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return value == pattern
	}

	position := 0
	for index, part := range parts {
		if part == "" {
			continue
		}

		foundAt := strings.Index(value[position:], part)
		if foundAt < 0 {
			return false
		}

		foundAt += position
		if index == 0 && !strings.HasPrefix(value, part) {
			return false
		}

		position = foundAt + len(part)
	}

	lastPart := parts[len(parts)-1]
	if lastPart == "" {
		return true
	}

	return strings.HasSuffix(value, lastPart)
}
