package handler

import (
	"42tokyo-road-to-dojo-go/pkg/http/middleware"
	"encoding/json"
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

// HandleUserCreate ユーザー作成処理 nameからtokenを作成
func (handler *handler) HandleUserCreate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var requestBody userCreateRequest //構造体変数宣言
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if requestBody.Name == "" {
			log.Println("name is empty")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		token := uuid.NewString()

		// データベースへのユーザー登録
		if err := handler.userDao.Create(requestBody.Name, token); err != nil {
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
		interfaceUserID := ctx.Value(middleware.UserIDKey)

		//Valueが返す値はinterface{}型のため変換
		userID, ok := interfaceUserID.(int) //okはbool
		if !ok {
			//500エラー
			log.Println("Failed to assert userID from context")
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
