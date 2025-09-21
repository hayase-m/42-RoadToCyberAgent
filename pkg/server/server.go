package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"42tokyo-road-to-dojo-go/pkg/http/middleware"
	"42tokyo-road-to-dojo-go/pkg/infra/dao"
	"42tokyo-road-to-dojo-go/pkg/server/handler"

	_ "github.com/go-sql-driver/mysql"
)

// Serve HTTPサーバを起動する
func Serve(addr string) {
	// parseTime=true を指定しないと、TIME 型のカラムを Go の time.Time にマッピングできない
	db, err := sql.Open("mysql", "root:ca-tech-dojo@tcp(127.0.0.1:3306)/game_api_db?parseTime=true")
	if err != nil {
		log.Fatalf("DB connection failed: %+v", err)
	}
	defer db.Close()

	// 依存性の注入 (DI)
	userDao := dao.NewUserDao(db)
	itemDao := dao.NewItemDao(db)

	//アイテムのマスターデータをキャッシュ
	items, err := itemDao.FindAll(context.Background())
	if err != nil {
		log.Fatalf("Failed to get all items: %v", err)
	}

	appHandler := handler.NewHandler(userDao, itemDao, items)
	appMiddleware := middleware.NewMiddleware(userDao)

	/* ===== URLマッピングを行う ===== */
	http.HandleFunc("/setting/get", get(appHandler.HandleSettingGet()))
	http.HandleFunc("/user/create", post(appHandler.HandleUserCreate()))
	http.HandleFunc("/user/get", get(appMiddleware.Authenticate(appHandler.HandleUserGet())))
	http.HandleFunc("/user/update", post(appMiddleware.Authenticate(appHandler.HandleUserUpdate())))

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}

// get GETリクエストを処理する
func get(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodGet)
}

// post POSTリクエストを処理する
func post(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodPost)
}

// httpMethod 指定したHTTPメソッドでAPIの処理を実行する
func httpMethod(apiFunc http.HandlerFunc, method string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// CORS対応
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

		// プリフライトリクエストは処理を通さない
		if request.Method == http.MethodOptions {
			return
		}
		// 指定のHTTPメソッドでない場合はエラー
		if request.Method != method {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method Not Allowed"))
			return
		}

		// 共通のレスポンスヘッダを設定
		writer.Header().Add("Content-Type", "application/json")
		apiFunc(writer, request)
	}
}
