package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/egayurcel990/go-uptime-monitor/internal/alert"
	"github.com/egayurcel990/go-uptime-monitor/internal/checker"
	"github.com/egayurcel990/go-uptime-monitor/internal/config"
	"github.com/egayurcel990/go-uptime-monitor/internal/handler"
	"github.com/egayurcel990/go-uptime-monitor/internal/metrics"
	"github.com/egayurcel990/go-uptime-monitor/internal/repository"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	db, err := repository.NewSQLite(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	repo := repository.New(db)
	metrics.Register()
	alerter := alert.New(cfg.WebhookURL)
	chk := checker.New(repo, alerter, cfg)

	go chk.Start()

	e := handler.NewRouter(repo, chk)
	log.Printf("Server running on :%s", cfg.Port)
	log.Fatal(e.Start(":" + cfg.Port))
}
