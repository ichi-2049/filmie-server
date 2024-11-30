package services

import (
	domain "github.com/ichi-2049/filmie-server/internal/domain/models"
	"github.com/ichi-2049/filmie-server/internal/domain/repositories"
)

type MovieService struct {
	movieRepo repositories.MovieRepository
}

func NewMovieService(movieRepo repositories.MovieRepository) *MovieService {
	return &MovieService{
		movieRepo: movieRepo,
	}
}

func (s *MovieService) GetAllMovies() ([]*domain.Movie, error) {
	return s.movieRepo.GetAllMovies()
}

func (s *MovieService) GetMovieConnection(first int, after *string) (*domain.MovieConnection, error) {
	return s.movieRepo.GetMovieConnection(first, after)
}
