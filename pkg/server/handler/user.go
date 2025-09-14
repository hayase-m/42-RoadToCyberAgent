package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type userCreateRequest struct {
	Name string `json:"name"`
}

type userCreateResponse struct {
	Token string `json:"token"`
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
