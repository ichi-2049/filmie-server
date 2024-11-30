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

// ページネーション結果を格納する構造体
type MovieConnection struct {
	Edges      []*MovieEdge // 映画情報とカーソル情報を含むエッジの配列
	PageInfo   *PageInfo    // ページング情報（次ページの有無、最後のカーソル）
	TotalCount int64        // 総件数
}

// 各映画エントリーの情報を格納する構造体
type MovieEdge struct {
	Cursor string // この映画エントリーの位置を示すカーソル
	Node   *Movie // 映画の実データ
}

// ページング情報を格納する構造体
type PageInfo struct {
	HasNextPage bool    // 次のページが存在するかどうか
	EndCursor   *string // 現在のページの最後のカーソル
}
