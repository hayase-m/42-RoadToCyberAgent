package handler

import (
	"42tokyo-road-to-dojo-go/pkg/http/middleware"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type userCreateRequest struct {
	Name string `json:"name"`
}

type userCreateResponse struct {
	Token string `json:"token"`
}

type userGetResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	HighScore int    `json:"highScore"`
	Coin      int    `json:"coin"`
}

type userUpdateRequest struct {
	Name string `json:"name"`
}

type collection struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int    `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}
type collectionListResponse struct {
	Collections []collection `json:"collections"`
}

type rankInfo struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
}

type rankingListResponse struct {
	Ranks []rankInfo `json:"ranks"`
}

type gameFinishRequest struct {
	Score int `json:"score"`
}

type gameFinishResponse struct {
	Coin int `json:"coin"`
}

// HandleUserCreate ユーザー作成処理 nameからtokenを作成
func (handler *handler) HandleUserCreate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		var requestBody userCreateRequest //構造体変数宣言
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			handler.handleError(writer, err, http.StatusBadRequest)
			return
		}

		err := handler.validateName(requestBody.Name)
		if err != nil {
			handler.handleError(writer, err, http.StatusBadRequest)
			return
		}

		token := uuid.NewString()

		// データベースへのユーザー登録
		if err := handler.userDao.Create(ctx, requestBody.Name, token); err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		response := userCreateResponse{
			Token: token,
		}
		data, err := json.Marshal(response)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		writer.Write(data)
	}
}

// HandleUserGet ユーザーget処理 idからuserを入手
func (handler *handler) HandleUserGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userID, err := handler.getUserIDFromContext(ctx)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		user, err := handler.userDao.FindByID(ctx, userID)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		response := userGetResponse{
			ID:        strconv.Itoa(user.ID),
			Name:      user.Name,
			HighScore: user.Highscore,
			Coin:      user.Coin,
		}

		data, err := json.Marshal(response)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		writer.Write(data)
	}
}

// HandleUserUpdate ユーザー名変更処理
func (handler *handler) HandleUserUpdate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userID, err := handler.getUserIDFromContext(ctx)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		var requestBody userUpdateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			handler.handleError(writer, err, http.StatusBadRequest)
			return
		}

		err = handler.validateName(requestBody.Name)
		if err != nil {
			handler.handleError(writer, err, http.StatusBadRequest)
			return
		}

		err = handler.userDao.UpdateName(ctx, userID, requestBody.Name)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (handler *handler) HandleCollectionList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userID, err := handler.getUserIDFromContext(ctx)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		userCollectionItemIDs, err := handler.userDao.GetUserCollectionItemIDs(ctx, userID)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		collections := collectionListResponse{
			Collections: make([]collection, 0),
		}

		userItemMap := make(map[int]bool)
		for _, itemID := range userCollectionItemIDs {
			userItemMap[itemID] = true
		}

		for _, item := range handler.items {
			hasItem := userItemMap[item.ID] //mapはkeyが存在しなければ型のゼロ値を返す（今回はfalse）

			c := collection{
				CollectionID: strconv.Itoa(item.ID),
				Name:         item.Name,
				Rarity:       item.Rarity,
				HasItem:      hasItem,
			}
			collections.Collections = append(collections.Collections, c)
		}

		data, err := json.Marshal(collections)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		writer.Write(data)
	}
}

func (handler *handler) HandleRankingList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		usersCount, err := handler.userDao.CountAllUsers(ctx)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}
		startStr := request.URL.Query().Get("start")
		start, err := strconv.Atoi(startStr)
		if err != nil {
			handler.handleError(writer, err, http.StatusBadRequest)
			return
		}
		if start < 1 || start > usersCount {
			log.Printf("Invalid value for 'start' parameter")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		rankList, err := handler.userDao.GetRankingList(ctx, start, 10)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		ranks := rankingListResponse{
			Ranks: make([]rankInfo, 0),
		}

		for i, rank := range rankList {
			r := rankInfo{
				UserID:   strconv.Itoa(rank.UserID),
				UserName: rank.UserName,
				Rank:     start + i,
				Score:    rank.Score,
			}
			ranks.Ranks = append(ranks.Ranks, r)
		}

		data, err := json.Marshal(ranks)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		writer.Write(data)
	}
}

func (handler *handler) HandleGameFinish() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		var requestBody gameFinishRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			handler.handleError(writer, err, http.StatusBadRequest)
			return
		}

		userID, err := handler.getUserIDFromContext(ctx)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		score := requestBody.Score
		if score < 0 {
			log.Printf("Invalid value for 'score' parameter")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		err = handler.userDao.UpdateHighscore(ctx, userID, score)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		coin := score * 100 //なんでもいい
		err = handler.userDao.AddCoins(ctx, userID, coin)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		c := gameFinishResponse{
			Coin: coin,
		}

		data, err := json.Marshal(c)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		writer.Write(data)
	}
}

func (handler *handler) getUserIDFromContext(ctx context.Context) (int, error) {
	interfaceUserID := ctx.Value(middleware.UserIDKey)
	if interfaceUserID == nil {
		return 0, errors.New("userID not found in context")
	}

	//Valueが返す値はinterface{}型のため変換
	userID, ok := interfaceUserID.(int) //okはbool
	if !ok {
		//500エラー
		return 0, errors.New("userID in context is not int")
	}

	return userID, nil
}

func (handler *handler) validateName(name string) error {
	if name == "" {
		return errors.New("name is empty")
	}
	if len(name) > 30 {
		return errors.New("name is too long")
	}
	return nil
}
