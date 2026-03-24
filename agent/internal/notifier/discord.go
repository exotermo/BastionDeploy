package notifier

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Discord struct {
	webhookURL string
}

func NewDiscord(webhookURL string) *Discord {
	return &Discord{webhookURL: webhookURL}
}

func (d *Discord) Send(message string) {
	if d.webhookURL == "" {
		log.Println("⚠️  DISCORD_WEBHOOK_URL não configurado, pulando notificação")
		return
	}

	payload, _ := json.Marshal(map[string]string{
		"content": message,
	})

	resp, err := http.Post(d.webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("❌ Erro ao notificar Discord: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("📣 Discord notificado: %d", resp.StatusCode)
}