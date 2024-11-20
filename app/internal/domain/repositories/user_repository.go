package repositories

import domain "github.com/ichi-2049/filmie-server/internal/domain/models"

type UserRepository interface {
	GetUser(uid string) (*domain.User, error)
}
