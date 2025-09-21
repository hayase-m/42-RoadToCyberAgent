package dao

import (
	"database/sql"
)

// UserDaoという型を定義。
type UserDao interface {
	Create(name, token string) error
}

type userDao struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) UserDao {
	return &userDao{db: db}
}

// ユーザーを新規作成
func (userDao *userDao) Create(name, token string) error {
	_, err := userDao.db.Exec("INSERT INTO users (name, token) VALUES (?, ?)", name, token)
	if err != nil {
		return err
	}

	return nil
}
