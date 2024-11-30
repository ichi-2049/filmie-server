package repositories

import domain "github.com/ichi-2049/filmie-server/internal/domain/models"

type MovieRepository interface {
	GetAllMovies() ([]*domain.Movie, error)
	GetMovieConnection(first int, after *string) (*domain.MovieConnection, error)
}
