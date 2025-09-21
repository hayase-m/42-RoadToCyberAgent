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

type collectionListRequest struct {
	Token string `json:"token"`
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

// HandleUserCreate ユーザー作成処理 nameからtokenを作成
func (handler *handler) HandleUserCreate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		var requestBody userCreateRequest //構造体変数宣言
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		err := handler.validateName(requestBody.Name)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		token := uuid.NewString()

		// データベースへのユーザー登録
		if err := handler.userDao.Create(ctx, requestBody.Name, token); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := userCreateResponse{
			Token: token,
		}
		data, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
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
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := handler.userDao.FindByID(ctx, userID)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
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
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
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
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var requestBody userUpdateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		err = handler.validateName(requestBody.Name)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		err = handler.userDao.UpdateName(ctx, userID, requestBody.Name)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
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
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		userCollectionItemIDs, err := handler.userDao.GetUserCollectionItemIDs(ctx, userID)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
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
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
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
