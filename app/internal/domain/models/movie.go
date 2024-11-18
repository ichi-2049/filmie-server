package domain

import (
	"github.com/ichi-2049/filmie-server/types"
)

type Movie struct {
	MovieID          int        `json:"movieId"`
	Title            string     `json:"title"`
	Overview         string     `json:"overview"`
	ReleaseDate      types.Date `json:"releaseDate"`
	TMDBImageURL     string     `json:"tmdbImageUrl"`
	S3ImageURL       string     `json:"s3ImageUrl"`
	Popularity       float32    `json:"popularity"`
	OriginalLanguage string     `json:"originalLanguage"`
	VoteAverage      float32    `json:"voteAverage"`
	VoteCount        uint32     `json:"voteCount"`
}
