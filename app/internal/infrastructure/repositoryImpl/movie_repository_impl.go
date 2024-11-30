package repositoryImpl

import (
	"encoding/base64"
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

// 疎通テスト用関数　作品情報を取得する関数
func (r *MovieRepositoryImpl) GetAllMovies() ([]*domain.Movie, error) {
	var movieDao []*dao.Movie
	result := r.db.Find(&movieDao).Limit(10)
	if result.Error != nil {
		return nil, result.Error
	}

	movies := make([]*domain.Movie, len(movieDao))
	for i, dao := range movieDao {
		movies[i] = dao.ToModel()
	}

	return movies, nil
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

// カーソルベースのページネーションを使用して映画情報を取得する関数
// first: 取得する件数
// after: 前回の最後のカーソル（次のページの開始位置）
func (r *MovieRepositoryImpl) GetMovieConnection(first int, after *string) (*domain.MovieConnection, error) {
	// クエリのベースを作成
	query := r.db.Model(&dao.Movie{})

	// 総件数を取得
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// afterカーソルが指定されている場合、その位置以降のデータを取得
	if after != nil {
		decodedCursor, err := decodeCursor(*after)
		if err != nil {
			return nil, err
		}
		query = query.Where("movie_id > ?", decodedCursor)
	}

	// 次ページの有無を確認するため、要求された件数+1を取得
	limit := first + 1
	var movieDao []*dao.Movie
	if err := query.Order("popularity DESC").Limit(limit).Find(&movieDao).Error; err != nil {
		return nil, err
	}

	// 次ページの有無を判定
	// 要求された件数より多くデータが取得できた場合、次ページが存在する
	hasNextPage := len(movieDao) > first
	if hasNextPage {
		// 次ページ判定用の余分なデータを削除
		movieDao = movieDao[:first]
	}

	// MovieEdgeの配列を作成
	edges := make([]*domain.MovieEdge, len(movieDao))
	for i, movie := range movieDao {
		edges[i] = &domain.MovieEdge{
			Cursor: encodeCursor(movie.MovieID), // movieIDをbase64エンコードしてカーソルを生成
			Node:   movie.ToModel(),             // DAOをドメインモデルに変換
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

// movieIDをカーソル文字列にエンコードする関数
func encodeCursor(movieID int) string {
	// "movie:{id}" の形式でbase64エンコード
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("movie:%d", movieID)))
}

// カーソル文字列をmovieIDにデコードする関数
func decodeCursor(cursor string) (int, error) {
	// base64デコード
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, err
	}

	// "movie:{id}" 形式から数値を抽出
	var movieID int
	_, err = fmt.Sscanf(string(decoded), "movie:%d", &movieID)
	if err != nil {
		return 0, err
	}

	return movieID, nil
}
