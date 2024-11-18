package dao

import (
	"time"

	domain "github.com/ichi-2049/filmie-server/internal/domain/models"
	"github.com/ichi-2049/filmie-server/types"
)

type Movie struct {
	MovieID          int       `gorm:"column:movie_id"`
	Title            string    `gorm:"column:title"`
	Overview         string    `gorm:"column:overview"`
	ReleaseDate      time.Time `gorm:"column:release_date"`
	TMDBImageURL     string    `gorm:"column:tmdb_image_url"`
	S3ImageURL       string    `gorm:"column:s3_image_url"`
	Popularity       float32   `gorm:"column:popularity"`
	OriginalLanguage string    `gorm:"column:original_language"`
	VoteAverage      float32   `gorm:"column:vote_average"`
	VoteCount        uint32    `gorm:"column:vote_count"`
}

func (d *Movie) ToModel() *domain.Movie {
	return &domain.Movie{
		MovieID:          d.MovieID,
		Title:            d.Title,
		Overview:         d.Overview,
		ReleaseDate:      *types.NewDate(d.ReleaseDate),
		TMDBImageURL:     d.TMDBImageURL,
		S3ImageURL:       d.S3ImageURL,
		Popularity:       d.Popularity,
		OriginalLanguage: d.OriginalLanguage,
		VoteAverage:      d.VoteAverage,
		VoteCount:        d.VoteCount,
	}
}

func FromModel(m *domain.Movie) *Movie {
	return &Movie{
		MovieID:          m.MovieID,
		Title:            m.Title,
		Overview:         m.Overview,
		ReleaseDate:      m.ReleaseDate.Time,
		TMDBImageURL:     m.TMDBImageURL,
		S3ImageURL:       m.S3ImageURL,
		Popularity:       m.Popularity,
		OriginalLanguage: m.OriginalLanguage,
		VoteAverage:      m.VoteAverage,
		VoteCount:        m.VoteCount,
	}
}
