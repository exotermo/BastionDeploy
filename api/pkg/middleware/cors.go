// regras de requisições

package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewCORS recebe a lista de origens permitidas e devolve o middleware configurado.
// middleware em Go é uma função que roda antes do handler, basicamente serve como um filto de café pros mais chegados
// igual interceptors no spring boot em java  

func NewCORS(allowedOrigins []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: allowedOrigins,

		// metodos que a api aceita
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},

		// headers que o frontend pode enviar
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"X-Webhook-Signature", // header customizado pro HMAC
		},

		// quanto tempo o browser pode salvar em cache as regras de CORS
		MaxAge: 12 * time.Hour, // 12 horas 

		/*
		acho um tempo aceitavel por não expor necessariamente a WAN, o frontend automatiza o processo de deploy 
		utilizando e abusando do tunelamento da cloudflared 
		*/
	})
}