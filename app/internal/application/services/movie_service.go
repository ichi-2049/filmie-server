package services

import (
	domain "github.com/ichi-2049/filmie-server/internal/domain/models"
	"github.com/ichi-2049/filmie-server/internal/domain/repositories"
)

const (
	defaultFirst = 20
	maxFirst     = 100
)

type MovieService struct {
	movieRepo repositories.MovieRepository
}

func NewMovieService(movieRepo repositories.MovieRepository) *MovieService {
	return &MovieService{
		movieRepo: movieRepo,
	}
}

/*
映画作品情報をカーソルページネーションで取得する関数
*/
func (s *MovieService) GetMovieConnection(first int, after *string, title *string) (*domain.MovieConnection, error) {
	// 最大取得件数を制限
	if first <= 0 {
		first = defaultFirst
	}
	if first > maxFirst {
		first = maxFirst
	}
	return s.movieRepo.GetMovieConnection(first, after, title)
}
