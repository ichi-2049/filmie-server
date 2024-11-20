package dao

import (
	domain "github.com/ichi-2049/filmie-server/internal/domain/models"
)

type User struct {
	Uid   string `gorm:"column:uid"`
	Name  string `gorm:"column:name"`
	Email string `gorm:"column:email"`
}

func (d *User) ToModel() *domain.User {
	return &domain.User{
		Uid:   d.Uid,
		Name:  d.Name,
		Email: d.Email,
	}
}

func (d *User) ToDao(m *domain.User) *User {
	return &User{
		Uid:   m.Uid,
		Name:  m.Name,
		Email: m.Email,
	}
}
