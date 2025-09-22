package handler

import (
	"42tokyo-road-to-dojo-go/pkg/infra/dao"
	"log"
	"net/http"
)

// handler ハンドラ構造体
type handler struct {
	userDao  dao.UserDao
	itemDao  dao.ItemDao
	gachaDao dao.GachaDao
	items    []*dao.Item
	itemMap  map[int]*dao.Item
	gachas   []*dao.Gacha
}

// 新しいhandlerを生成
func NewHandler(userDao dao.UserDao, itemDao dao.ItemDao, gachaDao dao.GachaDao, items []*dao.Item, itemMap map[int]*dao.Item, gachas []*dao.Gacha) *handler {
	return &handler{
		userDao:  userDao,
		itemDao:  itemDao,
		gachaDao: gachaDao,
		items:    items,
		itemMap:  itemMap,
		gachas:   gachas,
	}
}

func (h *handler) handleError(writer http.ResponseWriter, err error, statusCode int) {
	log.Println(err)
	writer.WriteHeader(statusCode)
}
