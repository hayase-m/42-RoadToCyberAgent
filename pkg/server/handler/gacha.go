package handler

import (
	"42tokyo-road-to-dojo-go/pkg/infra/dao"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
)

type gachaDrawRequest struct {
	Times int `json:"times"`
}

type gachaResult struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int    `json:"rarity"`
	IsNew        bool   `json:"isNew"`
}

type gachaDrawResponse struct {
	Results []gachaResult `json:"results"`
}

func (handler *handler) HandleGachaDraw() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		var requestBody gachaDrawRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			handler.handleError(writer, err, http.StatusBadRequest)
			return
		}

		userID, err := handler.getUserIDFromContext(ctx)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		times := requestBody.Times
		if times <= 0 || times > 100 {
			handler.handleError(writer, errors.New("invalid value for 'times' parameter"), http.StatusBadRequest)
			return
		}

		user, err := handler.userDao.FindByID(ctx, userID)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		//コインが足りないときの処理
		coin := times * GachaCoinConsumption
		if coin > user.Coin {
			handler.handleError(writer, errors.New("not enough coins"), http.StatusBadRequest)
			return
		}

		results := make([]*dao.Gacha, 0)
		for i := 0; i < times; i++ {
			result := handler.drawGacha(handler.gachas)
			results = append(results, result)
		}

		userCollectionItemIDs, err := handler.userDao.GetUserCollectionItemIDs(ctx, userID)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}
		userItemMap := make(map[int]bool)
		for _, itemID := range userCollectionItemIDs {
			userItemMap[itemID] = true
		}

		unownedItemIDs := make([]int, 0)

		for _, result := range results {
			if !userItemMap[result.ItemID] {
				unownedItemIDs = append(unownedItemIDs, result.ItemID)
			}
		}

		err = handler.userDao.ExecuteGachaDrawTransaction(ctx, userID, coin, unownedItemIDs)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		responseResults := gachaDrawResponse{
			Results: make([]gachaResult, 0),
		}

		for _, result := range results {
			r := gachaResult{
				CollectionID: strconv.Itoa(result.ItemID),
				Name:         handler.itemMap[result.ItemID].Name,
				Rarity:       handler.itemMap[result.ItemID].Rarity,
				IsNew:        !userItemMap[result.ItemID],
			}
			responseResults.Results = append(responseResults.Results, r)
		}

		data, err := json.Marshal(responseResults)
		if err != nil {
			handler.handleError(writer, err, http.StatusInternalServerError)
			return
		}

		writer.Write(data)
	}
}

func (handler *handler) drawGacha(gachas []*dao.Gacha) *dao.Gacha {
	totalWeight := 0
	for _, gacha := range gachas {
		totalWeight += gacha.Weight
	}

	// 0から (totalWeight - 1) までの範囲でランダムな整数を一つ生成する
	randomNumber := rand.Intn(totalWeight)

	for _, gacha := range gachas {
		randomNumber -= gacha.Weight
		if randomNumber < 0 {
			return gacha
		}
	}

	// 万一ループを抜けてしまっても、gachasが空じゃなければ最後のitemを返す
	if len(gachas) > 0 {
		return gachas[len(gachas)-1]
	}

	return nil
}
