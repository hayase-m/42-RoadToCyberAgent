package handler

import (
	"42tokyo-road-to-dojo-go/pkg/infra/dao"
	"log"
	"net/http"
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

func (h *handler) handleError(writer http.ResponseWriter, err error, statusCode int) {
	log.Println(err)
	writer.WriteHeader(statusCode)
}
