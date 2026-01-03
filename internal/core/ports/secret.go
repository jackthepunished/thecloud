package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
)

type SecretRepository interface {
	Create(ctx context.Context, secret *domain.Secret) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Secret, error)
	GetByName(ctx context.Context, name string) (*domain.Secret, error)
	List(ctx context.Context) ([]*domain.Secret, error)
	Update(ctx context.Context, secret *domain.Secret) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SecretService interface {
	CreateSecret(ctx context.Context, name, value, description string) (*domain.Secret, error)
	GetSecret(ctx context.Context, id uuid.UUID) (*domain.Secret, error)
	GetSecretByName(ctx context.Context, name string) (*domain.Secret, error)
	ListSecrets(ctx context.Context) ([]*domain.Secret, error)
	DeleteSecret(ctx context.Context, id uuid.UUID) error
}
