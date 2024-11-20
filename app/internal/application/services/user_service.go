package services

import (
	domain "github.com/ichi-2049/filmie-server/internal/domain/models"
	"github.com/ichi-2049/filmie-server/internal/domain/repositories"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUser(uid string) (*domain.User, error) {
	return s.userRepo.GetUser(uid)
}
