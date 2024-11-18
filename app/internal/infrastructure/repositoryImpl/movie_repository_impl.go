package repositoryImpl

import (
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
