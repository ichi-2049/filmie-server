package repositories

import domain "github.com/ichi-2049/filmie-server/internal/domain/models"

type MovieRepository interface {
	GetMovieConnection(first int, after *string, title *string) (*domain.MovieConnection, error)
}
