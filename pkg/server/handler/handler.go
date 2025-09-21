package handler

import (
	"42tokyo-road-to-dojo-go/pkg/infra/dao"
)

// handler ハンドラ構造体
type handler struct {
	userDao dao.UserDao
	itemDao dao.ItemDao
	items   []*dao.Item
}

// 新しいhandlerを生成
func NewHandler(userDao dao.UserDao, itemDao dao.ItemDao, items []*dao.Item) *handler {
	return &handler{
		userDao: userDao,
		itemDao: itemDao,
		items:   items,
	}
}
