package repositoryImpl

import (
	domain "github.com/ichi-2049/filmie-server/internal/domain/models"
	"github.com/ichi-2049/filmie-server/internal/infrastructure/dao"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

// userを取得する関数
func (r *UserRepositoryImpl) GetUser(uid string) (*domain.User, error) {
	var userDao *dao.User
	result := r.db.First(&userDao, "uid = ?", uid)
	if result.Error != nil {
		return nil, result.Error
	}
	user := userDao.ToModel()

	return user, nil
}
