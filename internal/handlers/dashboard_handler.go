// Package httphandlers provides HTTP handlers for the API.
package httphandlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/pkg/httputil"
)

// DashboardHandler handles dashboard API endpoints.
type DashboardHandler struct {
	svc            ports.DashboardService
	allowedOrigins []string
}

// NewDashboardHandler creates a new dashboard handler.
//
// allowedOrigins is the explicit allowlist of cross-origin Origin values
// permitted to subscribe to the SSE stream. Accepts either repeated entries
// or a single comma-separated string (matching the WS_ALLOWED_ORIGINS form).
// An empty list means same-origin only — no Access-Control-Allow-Origin
// header is emitted, so browsers fall back to same-origin enforcement. The
// literal "*" entry opts into permissive mode for non-browser clients; even
// then the response echoes the request Origin rather than "*", because the
// SSE EventSource API ignores wildcards when credentials are involved. See #347.
func NewDashboardHandler(svc ports.DashboardService, allowedOrigins ...string) *DashboardHandler {
	cleaned := make([]string, 0, len(allowedOrigins))
	for _, raw := range allowedOrigins {
		for _, o := range strings.Split(raw, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}
	}
	return &DashboardHandler{svc: svc, allowedOrigins: cleaned}
}

// GetSummary returns resource counts and overview metrics.
// GET /api/dashboard/summary
// GetSummary returns resource counts and overview metrics
// @Summary Get dashboard summary
// @Description Gets counts for instances, volumes, vpcs, and events
// @Tags dashboard
// @Produce json
// @Security APIKeyAuth
// @Success 200 {object} domain.ResourceSummary
// @Router /api/dashboard/summary [get]
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	summary, err := h.svc.GetSummary(c.Request.Context())
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.Success(c, http.StatusOK, summary)
}

// GetRecentEvents returns the most recent audit events.
// GET /api/dashboard/events?limit=10
// GetRecentEvents returns the most recent audit events
// @Summary Get recent events
// @Description Gets a list of the latest audit/status events
// @Tags dashboard
// @Produce json
// @Security APIKeyAuth
// @Param limit query int false "Number of events to return" default(10)
// @Success 200 {array} domain.Event
// @Router /api/dashboard/events [get]
func (h *DashboardHandler) GetRecentEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Cap at 100 events
	}

	events, err := h.svc.GetRecentEvents(c.Request.Context(), limit)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.Success(c, http.StatusOK, events)
}

// GetStats returns full dashboard statistics.
// GET /api/dashboard/stats
// GetStats returns full dashboard statistics
// @Summary Get detailed dashboard stats
// @Description Gets comprehensive counts and aggregated resource information
// @Tags dashboard
// @Produce json
// @Security APIKeyAuth
// @Success 200 {object} domain.DashboardStats
// @Router /api/dashboard/stats [get]
func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats(c.Request.Context())
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.Success(c, http.StatusOK, stats)
}

// StreamEvents sends real-time dashboard updates via SSE.
// GET /api/dashboard/stream
// StreamEvents sends real-time dashboard updates via SSE
// @Summary Stream dashboard updates (SSE)
// @Description Real-time stream of dashboard summary updates via Server-Sent Events
// @Tags dashboard
// @Produce text/event-stream
// @Security APIKeyAuth
// @Router /api/dashboard/stream [get]
func (h *DashboardHandler) StreamEvents(c *gin.Context) {
	// Enforce CORS before any streaming headers are written. SSE inherits the
	// caller's cookies/API key, so accepting "*" would let an attacker-hosted
	// page receive a logged-in user's events.
	if !h.applyCORS(c) {
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Send initial summary
	summary, err := h.svc.GetSummary(c.Request.Context())
	if err == nil {
		c.SSEvent("summary", summary)
		c.Writer.Flush()
	}

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			summary, err := h.svc.GetSummary(c.Request.Context())
			if err != nil {
				continue
			}
			c.SSEvent("summary", summary)
			c.Writer.Flush()
		}
	}
}

// applyCORS validates the request Origin and, if accepted, writes the
// CORS response headers. Returns false (after writing 403) when the origin
// is set but not allowed; the caller must abort.
//
// This endpoint takes authoritative control of its CORS headers, clearing
// the permissive wildcard set by httputil.CORS() upstream so that an
// authenticated SSE stream can never be subscribed to by a third-party origin.
//
// Behaviour:
//   - Empty Origin → same-origin or non-browser request; clear wildcard
//     defaults and emit no Access-Control-* headers (browsers default to
//     same-origin).
//   - Origin in allowlist → echo it back with Vary: Origin and credentials.
//   - "*" in allowlist → echo the request Origin (never literal "*", which
//     EventSource rejects when credentials are involved).
//   - Otherwise → 403.
func (h *DashboardHandler) applyCORS(c *gin.Context) bool {
	headers := c.Writer.Header()
	headers.Del("Access-Control-Allow-Origin")
	headers.Del("Access-Control-Allow-Credentials")

	origin := c.GetHeader("Origin")
	if origin == "" {
		return true
	}

	allowed := false
	for _, o := range h.allowedOrigins {
		if o == "*" || o == origin {
			allowed = true
			break
		}
	}
	if !allowed {
		c.AbortWithStatus(http.StatusForbidden)
		return false
	}

	headers.Set("Access-Control-Allow-Origin", origin)
	headers.Set("Access-Control-Allow-Credentials", "true")
	headers.Set("Vary", "Origin")
	return true
}
