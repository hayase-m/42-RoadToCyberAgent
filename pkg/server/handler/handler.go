package handler

import (
	"42tokyo-road-to-dojo-go/pkg/infra/dao"
)

// handler ハンドラ構造体
type handler struct {
	userDao dao.UserDao
}

// NewHandler は新しいhandlerを生成します。
func NewHandler(userDao dao.UserDao) *handler {
	return &handler{
		userDao: userDao,
	}
}
