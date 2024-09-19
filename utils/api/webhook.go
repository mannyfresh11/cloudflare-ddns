package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var webhookURL = os.Getenv("WEBHOOKURL")

type WebhookData struct {
	Content string `json:"content"`
}

func SendHook(msg string) error {

	if webhookURL == "" {
		return nil
	}

	data := WebhookData{
		Content: msg,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Unabnle to marshal. %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Unabnle to send to discord. %v", err)
	}
	defer resp.Body.Close()

	return nil
}
