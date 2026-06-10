package routes

import (
	db "auth/internal/db/storage"
	"auth/internal/model"
	"auth/internal/utils/response"
	"encoding/json"
	"fmt"
	"net/http"
)

func Register(storage db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.RegisterRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

	}
}
