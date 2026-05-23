package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Alerter struct{ webhookURL string }

func New(webhookURL string) *Alerter { return &Alerter{webhookURL: webhookURL} }

// SendDown sends a webhook alert when a target goes down.
func (a *Alerter) SendDown(name, url, errMsg string) {
	a.send(fmt.Sprintf(":red_circle: *%s* is DOWN\nURL: %s\nError: %s", name, url, errMsg))
}

// SendUp sends a webhook alert when a target recovers.
func (a *Alerter) SendUp(name, url string) {
	a.send(fmt.Sprintf(":green_circle: *%s* is back UP\nURL: %s", name, url))
}

func (a *Alerter) send(text string) {
	if a.webhookURL == "" {
		return
	}
	body, _ := json.Marshal(map[string]string{"text": text})
	resp, err := http.Post(a.webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("alert webhook failed: %v", err)
		return
	}
	defer resp.Body.Close()
}
