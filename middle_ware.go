package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/VSM1le/malamala/internal/auth"
	"github.com/VSM1le/malamala/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithERROR(w, 403, fmt.Sprintf("Auth erro :%v", err))

		}
		user, err := apiCfg.DB.GetUserByApi(r.Context(), sql.NullString{String: apiKey})
		if err != nil {
			respondWithERROR(w, 400, fmt.Sprintf("Couldn't get user:%v", err))
		}
		handler(w, r, user)
	}
}
