package dao

import (
	"context"
	"database/sql"
)

// UserDaoという型を定義。
type UserDao interface {
	Create(ctx context.Context, name string, token string) error
	FindByToken(ctx context.Context, token string) (*User, error)
	FindByID(ctx context.Context, userID int) (*User, error)
	UpdateName(ctx context.Context, userID int, newName string) error
	GetUserCollectionItemIDs(ctx context.Context, userID int) ([]int, error)
	GetRankingList(ctx context.Context, start int, limit int) ([]*RankInfo, error)
	CountAllUsers(ctx context.Context) (int, error)
	UpdateHighscore(ctx context.Context, userID int, score int) error
	AddCoins(ctx context.Context, userID int, coin int) error
}

type userDao struct {
	db *sql.DB
}

// フィールド名大文字始まりで外部からアクセス可能。daoだけではなく、handlerからもアクセス可能にする
type User struct {
	ID        int
	Name      string
	Highscore int
	Coin      int
	Token     string
}

type RankInfo struct {
	UserID   int
	UserName string
	Score    int
}

func NewUserDao(db *sql.DB) UserDao {
	return &userDao{db: db}
}

// ユーザーを新規作成
func (userDao *userDao) Create(ctx context.Context, name string, token string) error {
	_, err := userDao.db.ExecContext(ctx, "INSERT INTO users (name, token) VALUES (?, ?)", name, token)
	if err != nil {
		return err
	}

	return nil
}

// Tokenからユーザーを検索し、返す
func (userDao *userDao) FindByToken(ctx context.Context, token string) (*User, error) {
	var user User
	err := userDao.db.QueryRowContext(ctx, "SELECT id, name, highscore, coin, token FROM users WHERE token = ?", token).Scan(&user.ID, &user.Name, &user.Highscore, &user.Coin, &user.Token) //Scanで書き込み
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// IDからユーザーを検索し、返す
func (userDao *userDao) FindByID(ctx context.Context, userID int) (*User, error) {
	var user User
	err := userDao.db.QueryRowContext(ctx, "SELECT id, name, highscore, coin, token FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Name, &user.Highscore, &user.Coin, &user.Token)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (userDao *userDao) UpdateName(ctx context.Context, userID int, newName string) error {
	_, err := userDao.db.ExecContext(ctx, "UPDATE users SET name=? WHERE id=? ", newName, userID)
	if err != nil {
		return err
	}

	return nil
}

func (userDao *userDao) GetUserCollectionItemIDs(ctx context.Context, userID int) ([]int, error) {
	rows, err := userDao.db.QueryContext(ctx, "SELECT item_id FROM user_collections WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	itemIDs := make([]int, 0)

	for rows.Next() {
		var itemID int
		err := rows.Scan(&itemID)
		if err != nil {
			return nil, err
		}
		itemIDs = append(itemIDs, itemID)
	}
	return itemIDs, nil
}

func (userDao *userDao) GetRankingList(ctx context.Context, start int, limit int) ([]*RankInfo, error) {
	rows, err := userDao.db.QueryContext(ctx, "SELECT id, name, highscore FROM users ORDER BY highscore DESC, id ASC LIMIT ? OFFSET ?", limit, start-1)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rankList := make([]*RankInfo, 0)

	for rows.Next() {
		var rankInfo RankInfo
		err := rows.Scan(&rankInfo.UserID, &rankInfo.UserName, &rankInfo.Score)
		if err != nil {
			return nil, err
		}
		rankList = append(rankList, &rankInfo)
	}
	return rankList, nil
}

func (userDao *userDao) CountAllUsers(ctx context.Context) (int, error) {
	var count int
	err := userDao.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (userDao *userDao) UpdateHighscore(ctx context.Context, userID int, score int) error {
	_, err := userDao.db.ExecContext(ctx, "UPDATE users SET highscore = ? WHERE highscore < ? AND id = ?", score, score, userID)
	if err != nil {
		return err
	}
	return nil
}

func (userDao *userDao) AddCoins(ctx context.Context, userID int, coin int) error {
	_, err := userDao.db.ExecContext(ctx, "UPDATE users SET coin = coin + ? WHERE id = ?", coin, userID)
	if err != nil {
		return err
	}
	return nil
}
