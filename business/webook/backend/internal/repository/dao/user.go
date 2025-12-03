package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        int64  `gorm:"primaryKey,autoIncrement"`
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt int64
	UpdatedAt int64
}

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now

	return dao.db.WithContext(ctx).Create(&user).Error
}
