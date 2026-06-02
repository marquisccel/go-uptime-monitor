package handler

import (
	"net/http"
	"strconv"

	"github.com/egayurcel990/go-uptime-monitor/internal/checker"
	"github.com/egayurcel990/go-uptime-monitor/internal/model"
	"github.com/egayurcel990/go-uptime-monitor/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	repo    *repository.Repository
	checker *checker.Checker
}

func NewRouter(repo *repository.Repository, chk *checker.Checker) *echo.Echo {
	h := &Handler{repo: repo, checker: chk}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.File("web/index.html")
	})

	e.GET("/healthz", h.HealthCheck)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	v1 := e.Group("/api/v1")
	v1.GET("/status", h.GetStatus)
	v1.GET("/targets", h.GetTargets)
	v1.POST("/targets", h.CreateTarget)
	v1.DELETE("/targets/:id", h.DeleteTarget)
	v1.GET("/targets/:id/history", h.GetHistory)
	v1.POST("/targets/:id/check", h.CheckTargetNow)

	return e
}

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) GetStatus(c echo.Context) error {
	summary, err := h.repo.GetUptimeSummary()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, summary)
}

func (h *Handler) GetTargets(c echo.Context) error {
	targets, err := h.repo.GetTargets()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, targets)
}

func (h *Handler) CreateTarget(c echo.Context) error {
	var t model.Target

	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if t.Interval == 0 {
		t.Interval = 60
	}

	if err := h.repo.CreateTarget(&t); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, t)
}

func (h *Handler) DeleteTarget(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	if err := h.repo.DeleteTarget(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetHistory(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	history, err := h.repo.GetHistory(id, 100)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, history)
}

func (h *Handler) CheckTargetNow(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	result, err := h.checker.CheckNow(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
