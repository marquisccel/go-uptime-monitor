package repository

import (
	"database/sql"
	"time"

	"github.com/egayurcel990/go-uptime-monitor/internal/model"
)

type Repository struct{ db *sql.DB }

func New(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) CreateTarget(t *model.Target) error {
	res, err := r.db.Exec(
		`INSERT INTO targets (name, url, interval) VALUES (?, ?, ?)`,
		t.Name, t.URL, t.Interval,
	)
	if err != nil {
		return err
	}

	t.ID, err = res.LastInsertId()
	t.CreatedAt = time.Now()
	return err
}

func (r *Repository) GetTargets() ([]model.Target, error) {
	rows, err := r.db.Query(`SELECT id, name, url, interval, created_at FROM targets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts []model.Target
	for rows.Next() {
		var t model.Target
		if err := rows.Scan(&t.ID, &t.Name, &t.URL, &t.Interval, &t.CreatedAt); err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}

	return ts, rows.Err()
}

func (r *Repository) GetTarget(id int64) (*model.Target, error) {
	var t model.Target
	err := r.db.QueryRow(
		`SELECT id, name, url, interval, created_at FROM targets WHERE id = ?`,
		id,
	).Scan(&t.ID, &t.Name, &t.URL, &t.Interval, &t.CreatedAt)

	return &t, err
}

func (r *Repository) DeleteTarget(id int64) error {
	_, err := r.db.Exec(`DELETE FROM targets WHERE id = ?`, id)
	return err
}

func (r *Repository) SaveCheckResult(res *model.CheckResult) error {
	row, err := r.db.Exec(
		`INSERT INTO check_results (target_id, status_code, latency_ms, is_up, error, checked_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		res.TargetID,
		res.StatusCode,
		res.Latency.Milliseconds(),
		res.IsUp,
		res.Error,
		res.CheckedAt,
	)
	if err != nil {
		return err
	}

	res.ID, err = row.LastInsertId()
	return err
}

func (r *Repository) GetHistory(targetID int64, limit int) ([]model.CheckResult, error) {
	rows, err := r.db.Query(
		`SELECT id, target_id, status_code, latency_ms, is_up, error, checked_at
		 FROM check_results
		 WHERE target_id = ?
		 ORDER BY checked_at DESC
		 LIMIT ?`,
		targetID,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.CheckResult
	for rows.Next() {
		var cr model.CheckResult
		var latMs int64

		if err := rows.Scan(
			&cr.ID,
			&cr.TargetID,
			&cr.StatusCode,
			&latMs,
			&cr.IsUp,
			&cr.Error,
			&cr.CheckedAt,
		); err != nil {
			return nil, err
		}

		cr.Latency = time.Duration(latMs) * time.Millisecond
		results = append(results, cr)
	}

	return results, rows.Err()
}

func (r *Repository) GetUptimeSummary() ([]model.UptimeSummary, error) {
	rows, err := r.db.Query(`
		SELECT
			t.id,
			t.name,
			t.url,

			COALESCE(
				(
					SELECT cr2.is_up
					FROM check_results cr2
					WHERE cr2.target_id = t.id
					ORDER BY cr2.checked_at DESC
					LIMIT 1
				),
				0
			) AS is_up,

			COALESCE(
				ROUND(
					AVG(
						CASE WHEN cr.is_up THEN 1.0 ELSE 0.0 END
					) * 100,
					2
				),
				0
			) AS uptime_pct,

			COALESCE(AVG(cr.latency_ms), 0) AS avg_latency,

			COALESCE(
				(
					SELECT cr3.checked_at
					FROM check_results cr3
					WHERE cr3.target_id = t.id
					ORDER BY cr3.checked_at DESC
					LIMIT 1
				),
				'never'
			) AS last_check

		FROM targets t
		LEFT JOIN check_results cr
			ON cr.target_id = t.id
			AND cr.checked_at >= datetime('now', '-24 hours')
		GROUP BY t.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []model.UptimeSummary
	for rows.Next() {
		var s model.UptimeSummary

		if err := rows.Scan(
			&s.TargetID,
			&s.Name,
			&s.URL,
			&s.IsUp,
			&s.UptimePct,
			&s.AvgLatency,
			&s.LastCheck,
		); err != nil {
			return nil, err
		}

		summaries = append(summaries, s)
	}

	return summaries, rows.Err()
}
