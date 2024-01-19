package httphelper

import (
	"encoding/json"
	"net/http"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/userhelper"
)

func LoginHandler[T any, MODEL database.TableWithID[IDTYPE], IDTYPE int64 | string](h userhelper.UserHelper[T, MODEL, IDTYPE],
	conditionFactory ConditionFactory[T]) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			SendError("Invalid request body", err, resp)
			return
		}

		condition, err := conditionFactory(req)
		if err != nil {
			SendError("Invalid condition for "+h.GetTableName(), err, resp)
			return
		}

		token, err := h.Login(req.Context(), body.Username, body.Password, condition)
		if err != nil {
			SendError("Error while login", err, resp)
			return
		}

		SendData(map[string]string{"token": token}, resp)
	}
}

func UpdatePasswordWithTokenHandler[T any, MODEL database.TableWithID[IDTYPE], IDTYPE int64 | string](h userhelper.UserHelper[T, MODEL, IDTYPE],
	conditionFactory ConditionFactory[T], secret string) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		var body struct {
			Token    string `json:"token"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			SendError("Invalid request body", err, resp)
			return
		}

		condition, err := conditionFactory(req)
		if err != nil {
			SendError("Invalid condition for "+h.GetTableName(), err, resp)
			return
		}

		err = h.UpdatePasswordFromToken(req.Context(), body.Token, body.Password, condition)
		if err != nil {
			SendError("Error while updating password", err, resp)
			return
		}

		SendData(map[string]string{"message": "Password updated successfully"}, resp)
	}
}
