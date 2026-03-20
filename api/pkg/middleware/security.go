// headers de segurança 
// necessarios pro Cloudflare)


package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders adiciona headers HTTP que protegem contra ataques comuns.
// Especialmente importante quando exposto pelo Cloudflare.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {

		// impede que o navegador xereta adivinhe o tipo do conteúdo
		c.Header("X-Content-Type-Options", "nosniff")

		// bloqueia a API de ser carregada dentro de um <iframe>
		c.Header("X-Frame-Options", "DENY")

		// força HTTPS, o Cloudflare já faz isso, mas é uma segunda camada
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// política de segurança 
		// só aceita recursos da própria origem
		c.Header("Content-Security-Policy", "default-src 'self'")

		// remove o header que revela que é um servidor Gin/Go
		c.Header("Server", "")

		// Sucesso
		c.Next()
	}
}