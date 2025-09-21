package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"42tokyo-road-to-dojo-go/pkg/infra/dao"
)

// ctxに直接strを渡すとエラーがでるので、新しい型を定義
type contextKey string

const UserIDKey contextKey = "user_id"

type middleware struct {
	userDao dao.UserDao
}

// 新しいmiddlewareを生成
func NewMiddleware(userDao dao.UserDao) *middleware {
	return &middleware{
		userDao: userDao,
	}
}

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func (middleware *middleware) Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		token := request.Header.Get("x-token")
		if token == "" {
			//認証エラーは401
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		user, err := middleware.userDao.FindByToken(ctx, token)
		if err != nil {
			if err == sql.ErrNoRows {
				//不正tokenは401
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			//それ以外は500エラー
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, UserIDKey, user.ID) //ctxはイミュータブルなので、新たなcontextを作成
		nextFunc(writer, request.WithContext(ctx))
	}
}
