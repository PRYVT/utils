package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AuthMiddleware struct {
	tokenManager *TokenManager
}

func NewAuthMiddleware(tm *TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenManager: tm}
}

func GetTokenFromHeader(c *gin.Context) string {

	// Retrieve the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Error().Msg("Authorization header is missing")
		return ""
	}

	// Check if the token has the Bearer prefix
	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		log.Error().Msg("Authorization header format must be Bearer {token}")
		return ""
	}

	// Extract the token
	tokenString := authHeader[len(prefix):]
	return tokenString
}

func (am *AuthMiddleware) AuthenticateMiddleware(c *gin.Context) {

	tokenString := GetTokenFromHeader(c)
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Verify the token
	_, err := am.tokenManager.VerifyToken(tokenString)
	if err != nil {
		log.Error().Err(err).Msg("Token verification failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}
