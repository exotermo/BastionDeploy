package middleware

import (
	"crypto/hmac"   // biblioteca padrão do Go para HMAC
	"crypto/sha256" // algoritmo SHA256 o mesmo que o GitHub usa para encriptação
	"encoding/hex"  // converte bytes para string hexadecimal
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ValidateGitHubWebhook retorna um middleware (filtro de cafe) que bloqueia requisições
// cuja assinatura HMAC não bater com o secret configurado no .env
//
// Fluxo:
//   1. parser no header X-Hub-Signature-256 enviado pelo GitHub
//   2. parser no body da requisição
//   3. recalcula o HMAC usando a chave secret
//   4. compara com o que o GitHub enviou
//   5. se não bater -> 401 (nao autorizo not not kkkkkkkk), se bater com o secret (malandrao da sala vip) -> foward no processo

func ValidateGitHubWebhook(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// da um parser no header de assinatura
		// formato de encripação do gitHub: "sha256=abc123def456..."
		signature := c.GetHeader("X-Hub-Signature-256")
		if signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "header X-Hub-Signature-256 ausente",
			})
			return
		}

		// Lê o corpo da requisição em memória runtime (em tempo de execução)
		// precisamos do body pra calcular o HMAC
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "erro ao ler o corpo da requisição",
			})
			return
		}

		// IMPORTANTE: o Gin consome o body ao ler.
		// Precisamos "recolocar" os bytes para o handler conseguir ler depois.
		c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

		// calcula o HMAC-SHA256 do body usando nosso secret
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(bodyBytes)
		expectedMAC := "sha256=" + hex.EncodeToString(mac.Sum(nil))


		// compara as assinaturas
		// hmac.Equal faz comparação em tempo constante evita timing attacks
		// (ataques que adivinham o secret medindo o tempo de resposta)
		if !hmac.Equal([]byte(expectedMAC), []byte(signature)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "assinatura inválida — requisição rejeitada",
			})
			return
		}

		// sucesso! passa pro handler
		c.Next()
	}
}