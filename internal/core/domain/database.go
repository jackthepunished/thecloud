package domain

import (
	"time"

	"github.com/google/uuid"
)

type DatabaseEngine string

const (
	EnginePostgres DatabaseEngine = "postgres"
	EngineMySQL    DatabaseEngine = "mysql"
)

type DatabaseStatus string

const (
	DatabaseStatusCreating DatabaseStatus = "CREATING"
	DatabaseStatusRunning  DatabaseStatus = "RUNNING"
	DatabaseStatusStopped  DatabaseStatus = "STOPPED"
	DatabaseStatusDeleting DatabaseStatus = "DELETING"
	DatabaseStatusFailed   DatabaseStatus = "FAILED"
)

type Database struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Engine      DatabaseEngine
	Version     string
	Status      DatabaseStatus
	VpcID       *uuid.UUID
	ContainerID string
	Port        int
	Username    string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
