package middleware

import (
	"context"
	"net/http"
)

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc {
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

		nextFunc(writer, request.WithContext(ctx))
	}
}
