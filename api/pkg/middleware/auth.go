// Middleware de autenticação por API Key.
//
// Como funciona:
//   - Todos os endpoints GET /api/v1/* exigem o header X-API-Key
//   - O webhook do GitHub é exceção (usa HMAC validation)
//   - Health check (/api/v1/health) é público para load balancers
package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// RequireAPIKey valida o header X-API-Key.
// Se a variável de ambiente API_KEY não estiver definida, o middleware é no-op
// (modo dev). Em produção, API_KEY é obrigatória.
func RequireAPIKey() gin.HandlerFunc {
	key := os.Getenv("API_KEY")

	return func(c * gin.Context) {
		// Modo dev: sem chave configurada, passa direto
		if key == "" {
			c.Next()
			return
		}

		provided := c.GetHeader("X-API-Key")
		if provided == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "header X-API-Key ausente",
			})
			return
		}

		// Comparação em tempo constante (previne timing attack)
		providedHash := sha256.Sum256([]byte(provided))
		keyHash := sha256.Sum256([]byte(key))

		if subtle.ConstantTimeCompare(providedHash[:], keyHash[:]) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "chave de API inválida",
			})
			return
		}

		c.Next()
	}
}

// PublicEndpoint wrapper que isola rotas sem auth (ex: health)
// Usado internamente em routes.go
func PublicEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
