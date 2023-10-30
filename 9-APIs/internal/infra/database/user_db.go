package database

import (
	"github.com/christianferraz/goexpert/9-APIs/internal/entity"
	"gorm.io/gorm"
)

type UserDB struct {
	DB *gorm.DB
}

func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{
		DB: db,
	}
}

func (db *UserDB) Create(user *entity.User) error {
	return db.DB.Create(user).Error
}

func (db *UserDB) FindByEmail(emaild string) (*entity.User, error) {
	var user entity.User
	if err := db.DB.Where("email =?", emaild).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
