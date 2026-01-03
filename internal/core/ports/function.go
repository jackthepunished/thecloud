package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
)

type FunctionRepository interface {
	Create(ctx context.Context, f *domain.Function) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Function, error)
	GetByName(ctx context.Context, userID uuid.UUID, name string) (*domain.Function, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.Function, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CreateInvocation(ctx context.Context, i *domain.Invocation) error
	GetInvocations(ctx context.Context, functionID uuid.UUID, limit int) ([]*domain.Invocation, error)
}

type FunctionService interface {
	CreateFunction(ctx context.Context, name, runtime, handler string, code []byte) (*domain.Function, error)
	GetFunction(ctx context.Context, id uuid.UUID) (*domain.Function, error)
	ListFunctions(ctx context.Context) ([]*domain.Function, error)
	DeleteFunction(ctx context.Context, id uuid.UUID) error
	InvokeFunction(ctx context.Context, id uuid.UUID, payload []byte, async bool) (*domain.Invocation, error)
	GetFunctionLogs(ctx context.Context, id uuid.UUID, limit int) ([]*domain.Invocation, error)
}
