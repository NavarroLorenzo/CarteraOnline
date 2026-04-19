package auth

import (
	"cenz/backend/internal/shared/httpjson"
	"strings"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "authenticated_user_id"

func RequireAuth(tokenManager *TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if !strings.HasPrefix(authHeader, "Bearer ") {
			httpjson.Error(c, 401, "unauthorized", "Autenticación requerida")
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		userID, err := tokenManager.Parse(token)
		if err != nil {
			httpjson.Error(c, 401, "invalid_token", "Token inválido o vencido")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, userID)
		c.Next()
	}
}

func GetAuthenticatedUserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get(ContextUserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)
	return userID, ok
}

func AbortIfUnauthenticated(c *gin.Context) (int64, bool) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		httpjson.Error(c, 401, "unauthorized", "Autenticación requerida")
		c.Abort()
		return 0, false
	}

	return userID, true
}
