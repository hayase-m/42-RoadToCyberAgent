package dao

import (
	"context"
	"database/sql"
)

type GachaDao interface {
	FindAll(ctx context.Context) ([]*Gacha, error)
}

type gachaDao struct {
	db *sql.DB
}

type Gacha struct {
	ID     int
	ItemID int
	Weight int
}

func NewGachaDao(db *sql.DB) GachaDao {
	return &gachaDao{db: db}
}

func (gachaDao *gachaDao) FindAll(ctx context.Context) ([]*Gacha, error) {
	rows, err := gachaDao.db.QueryContext(ctx, "SELECT id, item_id, weight FROM gachas")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gachas := make([]*Gacha, 0)

	for rows.Next() {
		var gacha Gacha
		err := rows.Scan(&gacha.ID, &gacha.ItemID, &gacha.Weight)
		if err != nil {
			return nil, err
		}
		gachas = append(gachas, &gacha)
	}
	return gachas, nil
}
