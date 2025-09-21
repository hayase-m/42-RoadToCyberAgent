package middleware

import (
	"context"
	"database/sql"
	"net/http"

	"42tokyo-road-to-dojo-go/pkg/infra/dao"
)

// ctxに直接strを渡すとエラーがでるので、新しい型を定義
type contextKey string

const userIDkey contextKey = "user_id"

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

		// TODO: implement here
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
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, userIDkey, user.ID) //ctxはイミュータブルなので、新たなcontextを作成
		nextFunc(writer, request.WithContext(ctx))
	}
}
