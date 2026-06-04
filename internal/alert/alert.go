package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const cooldownDuration = 5 * time.Minute

type Alerter struct {
	webhookURL string
	mu         sync.Mutex
	lastDown   map[string]time.Time // key = target name
}

func New(webhookURL string) *Alerter {
	return &Alerter{
		webhookURL: webhookURL,
		lastDown:   make(map[string]time.Time),
	}
}

// SendDown sends a webhook alert when a target goes down.
// Calls within the cooldown window for the same target are silently dropped.
func (a *Alerter) SendDown(name, url, errMsg string) {
	a.mu.Lock()
	last, ok := a.lastDown[name]
	if ok && time.Since(last) < cooldownDuration {
		a.mu.Unlock()
		return
	}
	a.lastDown[name] = time.Now()
	a.mu.Unlock()

	a.send(fmt.Sprintf(":red_circle: *%s* is DOWN\nURL: %s\nError: %s", name, url, errMsg))
}

// SendUp sends a webhook alert when a target recovers and resets the cooldown.
func (a *Alerter) SendUp(name, url string) {
	a.mu.Lock()
	delete(a.lastDown, name)
	a.mu.Unlock()

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
