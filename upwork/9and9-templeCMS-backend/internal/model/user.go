package model

import (
	"context"
	"errors"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/constant"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/userhelper"
)

type userValidator struct{}

func NewUserValidator() userhelper.UserValidator[User] {
	return &userValidator{}
}

func (uv *userValidator) GetUsername(user *User) string {
	return user.Username
}

func (uv *userValidator) GetPassword(user *User) string {
	return user.Password
}

func (uv *userValidator) SetPassword(user *User, password string) {
	user.Password = password
}

type User struct {
	database.TableID[int64]
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func UserIDFromCTX(ctx context.Context) (int64, error) {
	user, ok := ctx.Value(constant.CtxKey_User).(*User)
	if !ok {
		return 0, errors.New("user not found in context")
	}

	return user.ID, nil
}

func NewUserHelper(db database.Connection) database.CrudHelper[database.MysqlCondition, User, int64] {
	return database.NewBaseHelper(db, "user", func(t *User) map[string]interface{} {
		return map[string]interface{}{
			"id":       &t.ID,
			"username": &t.Username,
			"password": &t.Password,
		}
	})
}
