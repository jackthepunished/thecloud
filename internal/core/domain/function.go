package domain

import (
	"time"

	"github.com/google/uuid"
)

type Function struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Runtime   string    `json:"runtime"`
	Handler   string    `json:"handler"`
	CodePath  string    `json:"code_path"`
	Timeout   int       `json:"timeout"`
	MemoryMB  int       `json:"memory_mb"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Invocation struct {
	ID         uuid.UUID  `json:"id"`
	FunctionID uuid.UUID  `json:"function_id"`
	Status     string     `json:"status"`
	StartedAt  time.Time  `json:"started_at"`
	EndedAt    *time.Time `json:"ended_at,omitempty"`
	DurationMs int        `json:"duration_ms"`
	StatusCode int        `json:"status_code"`
	Logs       string     `json:"logs"`
}
