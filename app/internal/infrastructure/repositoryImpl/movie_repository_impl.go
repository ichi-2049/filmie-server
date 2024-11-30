package repositoryImpl

import (
	"encoding/base32"
	"fmt"

	domain "github.com/ichi-2049/filmie-server/internal/domain/models"
	"github.com/ichi-2049/filmie-server/internal/infrastructure/dao"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MovieRepositoryImpl struct {
	db *gorm.DB
}

func NewMovieRepositoryImpl(db *gorm.DB) *MovieRepositoryImpl {
	return &MovieRepositoryImpl{
		db: db,
	}
}

// カーソルの構造体を定義
type Cursor struct {
	Popularity float32
	MovieID    int
}

/*
映画作品情報をカーソルページネーションで取得する関数
タイトル（optional）で検索をかけ、人気順（降順）と映画ID（昇順）でソートする
*/
func (r *MovieRepositoryImpl) GetMovieConnection(first int, after *string, title *string) (*domain.MovieConnection, error) {
	// クエリのベースを作成
	query := r.db.Model(&dao.Movie{})

	// タイトル検索条件を追加
	if title != nil && *title != "" {
		query = query.Where("title LIKE ?", *title+"%")
	}

	// 総件数を取得
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// afterカーソルが指定されている場合、その位置以降のデータを取得
	var lastPopularity float32
	var lastMovieID int
	if after != nil {
		decodedCursor, err := decodeCursor(*after)
		if err != nil {
			return nil, err
		}
		lastPopularity = decodedCursor.Popularity
		lastMovieID = decodedCursor.MovieID

		// 人気度と映画IDで複合的にフィルタリング
		query = query.Where(
			"(popularity < ?) OR (popularity = ? AND movie_id > ?)",
			lastPopularity,
			lastPopularity,
			lastMovieID,
		)
	}

	// 次ページの有無を確認するため、要求された件数+1を取得
	limit := first + 1
	var movieDao []*dao.Movie
	if err := query.
		Order("popularity DESC, movie_id ASC"). // 人気度降順、映画ID昇順でソート
		Limit(limit).
		Find(&movieDao).Error; err != nil {
		return nil, err
	}

	// 次ページの有無を判定
	hasNextPage := len(movieDao) > first
	if hasNextPage {
		// 次ページ判定用の余分なデータを削除
		movieDao = movieDao[:first]
	}

	// MovieEdgeの配列を作成
	edges := make([]*domain.MovieEdge, len(movieDao))
	for i, movie := range movieDao {
		edges[i] = &domain.MovieEdge{
			Cursor: encodeCursor(movie.Popularity, movie.MovieID), // 人気度と映画IDを組み合わせたカーソル
			Node:   movie.ToModel(),
		}
	}

	// 最後のカーソルを取得
	var endCursor *string
	if len(edges) > 0 {
		cursor := edges[len(edges)-1].Cursor
		endCursor = &cursor
	}

	// 結果を格納して返却
	return &domain.MovieConnection{
		Edges: edges,
		PageInfo: &domain.PageInfo{
			HasNextPage: hasNextPage,
			EndCursor:   endCursor,
		},
		TotalCount: totalCount,
	}, nil
}

// 映画作品情報を一括で登録する関数
// 主キーが重複した場合は更新する
func (r *MovieRepositoryImpl) BulkInsertMovies(movies []*domain.Movie) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	var movieDaoList []*dao.Movie
	for _, movie := range movies {
		movieDaoList = append(movieDaoList, dao.FromModel(movie))
	}

	// ON DUPLICATE KEY UPDATE
	if err := tx.Clauses(clause.OnConflict{
		UpdateAll: true, // 主キーが重複した場合にすべてのフィールドを更新
	}).Create(movieDaoList).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// 人気度と映画IDをカーソル文字列にエンコードする関数
func encodeCursor(popularity float32, movieID int) string {
	// "popularity:{popularity}:movie:{id}" の形式でbase32エンコード
	return base32.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("popularity:%.6f:movie:%d", popularity, movieID)),
	)
}

// カーソル文字列を人気度と映画IDにデコードする関数
func decodeCursor(cursor string) (Cursor, error) {
	// base32デコード
	decoded, err := base32.StdEncoding.DecodeString(cursor)
	if err != nil {
		return Cursor{}, err
	}

	// "popularity:{popularity}:movie:{id}" 形式から値を抽出
	var popularity float32
	var movieID int
	_, err = fmt.Sscanf(
		string(decoded),
		"popularity:%f:movie:%d",
		&popularity,
		&movieID,
	)
	if err != nil {
		return Cursor{}, err
	}

	return Cursor{
		Popularity: popularity,
		MovieID:    movieID,
	}, nil
}
