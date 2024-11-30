package converter

import (
	gqlmodel "github.com/ichi-2049/filmie-server/graphql/models"
	domain "github.com/ichi-2049/filmie-server/internal/domain/models"
)

func ConvertMovieConnection(movieConnection *domain.MovieConnection) *gqlmodel.MovieConnection {
	if movieConnection == nil {
		return nil
	}

	edges := make([]*gqlmodel.MovieEdge, len(movieConnection.Edges))
	for i, e := range movieConnection.Edges {
		edges[i] = &gqlmodel.MovieEdge{
			Cursor: e.Cursor,
			Node: &gqlmodel.Movie{
				MovieID:          e.Node.MovieID,
				Title:            e.Node.Title,
				Overview:         e.Node.Overview,
				ReleaseDate:      e.Node.ReleaseDate,
				ImageURL:         e.Node.TMDBImageURL,
				Popularity:       float64(e.Node.Popularity),
				OriginalLanguage: e.Node.OriginalLanguage,
				VoteAverage:      float64(e.Node.VoteAverage),
				VoteCount:        int(e.Node.VoteCount),
			},
		}
	}

	var endCursor *string
	if movieConnection.PageInfo != nil && movieConnection.PageInfo.EndCursor != nil {
		endCursor = movieConnection.PageInfo.EndCursor
	}

	return &gqlmodel.MovieConnection{
		Edges: edges,
		PageInfo: &gqlmodel.PageInfo{
			HasNextPage: movieConnection.PageInfo != nil && movieConnection.PageInfo.HasNextPage,
			EndCursor:   endCursor,
		},
		TotalCount: int(movieConnection.TotalCount),
	}
}
