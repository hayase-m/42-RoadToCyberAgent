package dao

import (
	"context"
	"database/sql"
)

type ItemDao interface {
	FindAll(ctx context.Context) ([]*Item, error)
}

type itemDao struct {
	db *sql.DB
}

type Item struct {
	ID     int
	Name   string
	Rarity int
}

func NewItemDao(db *sql.DB) ItemDao {
	return &itemDao{db: db}
}

func (itemDao *itemDao) FindAll(ctx context.Context) ([]*Item, error) {
	rows, err := itemDao.db.QueryContext(ctx, "SELECT id, name, rarity FROM items")
	if err != nil {
		return nil, err
	}
	//deferで関数終了時に実行
	defer rows.Close()

	items := make([]*Item, 0)

	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.Rarity)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}
