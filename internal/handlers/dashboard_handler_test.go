package httphandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type dashboardServiceMock struct {
	mock.Mock
}

func (m *dashboardServiceMock) GetSummary(ctx context.Context) (*domain.ResourceSummary, error) {
	args := m.Called(ctx)
	// Helper for checking return value
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	r0, _ := args.Get(0).(*domain.ResourceSummary)
	return r0, args.Error(1)
}

func (m *dashboardServiceMock) GetRecentEvents(ctx context.Context, limit int) ([]*domain.Event, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	r0, _ := args.Get(0).([]*domain.Event)
	return r0, args.Error(1)
}

func (m *dashboardServiceMock) GetStats(ctx context.Context) (*domain.DashboardStats, error) {
	args := m.Called(ctx)
	// Check for nil return
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	r0, _ := args.Get(0).(*domain.DashboardStats)
	return r0, args.Error(1)
}

func setupDashboardHandlerTest(_ *testing.T) (*dashboardServiceMock, *DashboardHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(dashboardServiceMock)
	handler := NewDashboardHandler(mockSvc)
	r := gin.New()
	return mockSvc, handler, r
}

func TestDashboardHandlerGetSummary(t *testing.T) {
	t.Parallel()
	mockSvc, handler, r := setupDashboardHandlerTest(t)
	defer mockSvc.AssertExpectations(t)

	r.GET("/summary", handler.GetSummary)

	summary := &domain.ResourceSummary{TotalInstances: 5}
	mockSvc.On("GetSummary", mock.Anything).Return(summary, nil)

	req, err := http.NewRequest("GET", "/summary", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var wrapper struct {
		Data domain.ResourceSummary `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &wrapper))
	assert.Equal(t, 5, wrapper.Data.TotalInstances)
}

func TestDashboardHandlerGetRecentEvents(t *testing.T) {
	t.Parallel()
	mockSvc, handler, r := setupDashboardHandlerTest(t)
	defer mockSvc.AssertExpectations(t)

	r.GET("/events", handler.GetRecentEvents)

	events := []*domain.Event{{ID: uuid.New(), Action: "TEST"}}
	mockSvc.On("GetRecentEvents", mock.Anything, 10).Return(events, nil)

	req, err := http.NewRequest("GET", "/events?limit=10", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var wrapper struct {
		Data []*domain.Event `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &wrapper))
	assert.Len(t, wrapper.Data, 1)
}

func TestDashboardHandlerGetStats(t *testing.T) {
	t.Parallel()
	mockSvc, handler, r := setupDashboardHandlerTest(t)
	defer mockSvc.AssertExpectations(t)

	r.GET("/stats", handler.GetStats)

	stats := &domain.DashboardStats{
		CPUHistory: []domain.MetricPoint{{Value: 10.1}},
	}
	mockSvc.On("GetStats", mock.Anything).Return(stats, nil)

	req, err := http.NewRequest("GET", "/stats", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var wrapper struct {
		Data domain.DashboardStats `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &wrapper))
	assert.Len(t, wrapper.Data.CPUHistory, 1)
	assert.InDelta(t, 10.1, wrapper.Data.CPUHistory[0].Value, 0.01)
}

func TestDashboardHandlerStreamEvents(t *testing.T) {
	t.Parallel()
	mockSvc, handler, r := setupDashboardHandlerTest(t)
	defer mockSvc.AssertExpectations(t)

	r.GET("/stream", handler.StreamEvents)

	summary := &domain.ResourceSummary{TotalInstances: 10}
	mockSvc.On("GetSummary", mock.Anything).Return(summary, nil)

	req, err := http.NewRequest("GET", "/stream", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	req = req.WithContext(ctx)

	go r.ServeHTTP(w, req)

	time.Sleep(100 * time.Millisecond)
	cancel()

	assert.Contains(t, w.Header().Get("Content-Type"), "text/event-stream")
	assert.Contains(t, w.Body.String(), "event:summary")
	// Same-origin (no Origin header) → must not emit a wildcard CORS header. See #347.
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestDashboardHandlerStreamEventsCORS(t *testing.T) {
	t.Parallel()

	runStream := func(t *testing.T, h *DashboardHandler, origin string) *httptest.ResponseRecorder {
		t.Helper()
		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.GET("/stream", h.StreamEvents)

		req, err := http.NewRequest("GET", "/stream", nil)
		require.NoError(t, err)
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		ctx, cancel := context.WithCancel(context.Background())
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		done := make(chan struct{})
		go func() {
			r.ServeHTTP(w, req)
			close(done)
		}()
		time.Sleep(50 * time.Millisecond)
		cancel()
		<-done
		return w
	}

	t.Run("AllowedOriginIsEchoed", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(dashboardServiceMock)
		mockSvc.On("GetSummary", mock.Anything).Return(&domain.ResourceSummary{}, nil)
		h := NewDashboardHandler(mockSvc, "https://dash.example.com,https://other.example.com")

		w := runStream(t, h, "https://dash.example.com")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "https://dash.example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
		assert.Equal(t, "Origin", w.Header().Get("Vary"))
	})

	t.Run("DisallowedOriginRejected", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(dashboardServiceMock)
		h := NewDashboardHandler(mockSvc, "https://dash.example.com")

		w := runStream(t, h, "https://attacker.example.com")

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
		mockSvc.AssertNotCalled(t, "GetSummary", mock.Anything)
	})

	t.Run("EmptyAllowlistRejectsCrossOrigin", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(dashboardServiceMock)
		h := NewDashboardHandler(mockSvc)

		w := runStream(t, h, "https://attacker.example.com")

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
		mockSvc.AssertNotCalled(t, "GetSummary", mock.Anything)
	})

	t.Run("WildcardEchoesOriginNotStar", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(dashboardServiceMock)
		mockSvc.On("GetSummary", mock.Anything).Return(&domain.ResourceSummary{}, nil)
		h := NewDashboardHandler(mockSvc, "*")

		w := runStream(t, h, "https://anything.example.com")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "https://anything.example.com", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("SameOriginEmitsNoCORSHeader", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(dashboardServiceMock)
		mockSvc.On("GetSummary", mock.Anything).Return(&domain.ResourceSummary{}, nil)
		h := NewDashboardHandler(mockSvc, "https://dash.example.com")

		w := runStream(t, h, "")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestDashboardHandlerGetRecentEventsLimits(t *testing.T) {
	t.Parallel()
	mockSvc, handler, r := setupDashboardHandlerTest(t)
	r.GET("/events", handler.GetRecentEvents)

	t.Run("Default", func(t *testing.T) {
		mockSvc.On("GetRecentEvents", mock.Anything, 10).Return([]*domain.Event{}, nil).Once()
		req, _ := http.NewRequest("GET", "/events", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid", func(t *testing.T) {
		mockSvc.On("GetRecentEvents", mock.Anything, 10).Return([]*domain.Event{}, nil).Once()
		req, _ := http.NewRequest("GET", "/events?limit=abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Cap", func(t *testing.T) {
		mockSvc.On("GetRecentEvents", mock.Anything, 100).Return([]*domain.Event{}, nil).Once()
		req, _ := http.NewRequest("GET", "/events?limit=200", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDashboardHandlerErrors(t *testing.T) {
	t.Parallel()
	mockSvc, handler, r := setupDashboardHandlerTest(t)
	r.GET("/summary", handler.GetSummary)
	r.GET("/events", handler.GetRecentEvents)
	r.GET("/stats", handler.GetStats)

	t.Run("SummaryError", func(t *testing.T) {
		mockSvc.On("GetSummary", mock.Anything).Return(nil, assert.AnError)
		req, _ := http.NewRequest("GET", "/summary", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("EventsError", func(t *testing.T) {
		mockSvc.On("GetRecentEvents", mock.Anything, 10).Return(nil, assert.AnError)
		req, _ := http.NewRequest("GET", "/events", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("StatsError", func(t *testing.T) {
		mockSvc.On("GetStats", mock.Anything).Return(nil, assert.AnError)
		req, _ := http.NewRequest("GET", "/stats", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
