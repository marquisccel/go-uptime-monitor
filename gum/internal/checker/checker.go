package checker

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/egayurcel990/go-uptime-monitor/internal/alert"
	"github.com/egayurcel990/go-uptime-monitor/internal/config"
	"github.com/egayurcel990/go-uptime-monitor/internal/metrics"
	"github.com/egayurcel990/go-uptime-monitor/internal/model"
	"github.com/egayurcel990/go-uptime-monitor/internal/repository"
)

type Checker struct {
	repo    *repository.Repository
	alerter *alert.Alerter
	cfg     *config.Config
	client  *http.Client
}

func New(repo *repository.Repository, alerter *alert.Alerter, cfg *config.Config) *Checker {
	return &Checker{
		repo:    repo,
		alerter: alerter,
		cfg:     cfg,
		client:  &http.Client{Timeout: time.Duration(cfg.CheckTimeout) * time.Second},
	}
}

// Start runs the check loop — call in a goroutine.
func (c *Checker) Start() {
	ticker := time.NewTicker(time.Duration(c.cfg.CheckInterval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		c.checkAll()
	}
}

// CheckNow triggers an immediate check for a specific target.
func (c *Checker) CheckNow(targetID int64) (*model.CheckResult, error) {
	t, err := c.repo.GetTarget(targetID)
	if err != nil {
		return nil, err
	}
	return c.checkTarget(t)
}

func (c *Checker) checkAll() {
	targets, err := c.repo.GetTargets()
	if err != nil {
		log.Printf("checkAll: %v", err)
		return
	}
	for _, t := range targets {
		t := t
		go func() {
			if _, err := c.checkTarget(&t); err != nil {
				log.Printf("check failed %s: %v", t.URL, err)
			}
		}()
	}
}

func (c *Checker) checkTarget(t *model.Target) (*model.CheckResult, error) {
	start := time.Now()
	result := &model.CheckResult{TargetID: t.ID, CheckedAt: start}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cfg.CheckTimeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.URL, nil)
	if err != nil {
		result.IsUp = false
		result.Error = err.Error()
	} else {
		resp, err := c.client.Do(req)
		result.Latency = time.Since(start)
		if err != nil {
			result.IsUp = false
			result.Error = err.Error()
			c.alerter.SendDown(t.Name, t.URL, err.Error())
		} else {
			resp.Body.Close()
			result.StatusCode = resp.StatusCode
			result.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 400
			if !result.IsUp {
				c.alerter.SendDown(t.Name, t.URL, http.StatusText(resp.StatusCode))
			}
		}
	}

	// Update Prometheus metrics
	status := "up"
	if !result.IsUp {
		status = "down"
	}
	upVal := 0.0
	if result.IsUp {
		upVal = 1.0
	}
	metrics.CheckUp.WithLabelValues(t.Name, t.URL).Set(upVal)
	metrics.ChecksTotal.WithLabelValues(t.Name, t.URL, status).Inc()
	if result.Latency > 0 {
		metrics.CheckDuration.WithLabelValues(t.Name, t.URL).Observe(result.Latency.Seconds())
	}

	if err := c.repo.SaveCheckResult(result); err != nil {
		log.Printf("save result failed target %d: %v", t.ID, err)
	}
	return result, nil
}
