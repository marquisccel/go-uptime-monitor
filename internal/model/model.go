package model

import "time"

// Target is a URL to be monitored.
type Target struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Interval  int       `json:"interval"` // seconds
	CreatedAt time.Time `json:"created_at"`
}

// CheckResult is the outcome of a single uptime check.
type CheckResult struct {
	ID         int64         `json:"id"`
	TargetID   int64         `json:"target_id"`
	StatusCode int           `json:"status_code"`
	Latency    time.Duration `json:"latency_ms"`
	IsUp       bool          `json:"is_up"`
	CheckedAt  time.Time     `json:"checked_at"`
	Error      string        `json:"error,omitempty"`
}

// UptimeSummary is the aggregated status for a target.
type UptimeSummary struct {
	TargetID   int64   `json:"target_id"`
	Name       string  `json:"name"`
	URL        string  `json:"url"`
	IsUp       bool    `json:"is_up"`
	UptimePct  float64 `json:"uptime_pct"`
	AvgLatency float64 `json:"avg_latency_ms"`
	LastCheck  string  `json:"last_check"`
}
