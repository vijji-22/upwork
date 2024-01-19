package userhelper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
)

type UserValidator[Model any] interface {
	GetUsername(m *Model) string
	GetPassword(m *Model) string
	SetPassword(m *Model, password string)
}

type UserHelper[T any, UserModel database.TableWithID[IDTYPE], IDTYPE int64 | string] interface {
	database.CrudHelper[T, UserModel, IDTYPE]
	Login(ctx context.Context, username string, password string, condition database.Condition[T]) (string, error)
	WithLoginValidationFields(fields ...string) *userHelper[T, UserModel, IDTYPE]
	UpdatePasswordFromToken(ctx context.Context, token, password string, condition database.Condition[T]) error
}

type userHelper[T any, UserModel database.TableWithID[IDTYPE], IDTYPE int64 | string] struct {
	database.CrudHelper[T, UserModel, IDTYPE]
	secret string

	userValidator UserValidator[UserModel]

	loginValidationFields []string
}

func NewUserHelper[T any, UserModel database.TableWithID[IDTYPE], IDTYPE int64 | string](helper database.CrudHelper[T, UserModel, IDTYPE],
	secret string, userValidator UserValidator[UserModel]) UserHelper[T, UserModel, IDTYPE] {
	return &userHelper[T, UserModel, IDTYPE]{
		CrudHelper:    helper,
		secret:        secret,
		userValidator: userValidator,

		loginValidationFields: []string{"username"},
	}
}

func (u *userHelper[T, UserModel, IDTYPE]) WithLoginValidationFields(fields ...string) *userHelper[T, UserModel, IDTYPE] {
	u.loginValidationFields = fields
	return u
}

func (u *userHelper[T, UserModel, IDTYPE]) Create(ctx context.Context, user *UserModel) (*UserModel, error) {
	if u.userValidator.GetUsername(user) == "" {
		return nil, errors.New("username is required")
	}

	if u.userValidator.GetPassword(user) == "" {
		return nil, errors.New("password is required")
	}

	pass, err := encryptPassword(u.userValidator.GetPassword(user))
	if err != nil {
		return nil, err
	}

	u.userValidator.SetPassword(user, pass)

	user, err = u.CrudHelper.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	u.userValidator.SetPassword(user, "") // removing password as it should not be available in reponse
	return user, nil
}

func (u *userHelper[T, UserModel, IDTYPE]) Update(ctx context.Context, user *UserModel, project []string, condition database.Condition[T]) error {
	if len(project) == 0 || project[0] == "" || project[0] == "*" {
		return errors.New("all field updates are not allowed")
	}

	projectSet := utils.NewSetFromSlice(project)

	pass, username := u.userValidator.GetPassword(user), u.userValidator.GetUsername(user)

	if pass == "" && projectSet.Contains("password") {
		return errors.New("password is required")
	}

	if username == "" && projectSet.Contains("username") {
		return errors.New("username is required")
	}

	if pass != "" {
		pass, err := encryptPassword(pass)
		if err != nil {
			return err
		}

		u.userValidator.SetPassword(user, pass)
	}

	return u.CrudHelper.Update(ctx, user, project, condition)
}

func (u *userHelper[T, UserModel, IDTYPE]) Get(ctx context.Context, project []string, condition database.Condition[T]) ([]UserModel, error) {
	if len(project) == 0 || project[0] == "" || project[0] == "*" {
		return nil, errors.New("all field get are not allowed")
	}

	projectSet := utils.NewSetFromSlice(project)
	projectSet.Remove("password")

	return u.CrudHelper.Get(ctx, project, condition)
}

func (u *userHelper[T, UserModel, IDTYPE]) Login(ctx context.Context, username string, password string, condition database.Condition[T]) (string, error) {
	user, err := u.getUserByUsername(ctx, username, condition)
	if err != nil {
		return "", err
	}

	if password == "" {
		return "", errors.New("password is required")
	}

	if !comparePassword(password, u.userValidator.GetPassword(user)) {
		return "", errors.New("invalid password")
	}

	u.userValidator.SetPassword(user, "") // removing this as it should not be available in token

	return GenerateJWT(user, "login", time.Hour, u.secret)
}

func (u *userHelper[T, UserModel, IDTYPE]) GenerateForgetPasswordToken(ctx context.Context, username string, condition database.Condition[T]) (string, error) {
	user, err := u.getUserByUsername(ctx, username, condition)
	if err != nil {
		return "", err
	}

	return GenerateJWT(user, "password_reset", time.Hour, u.secret)
}

func (u *userHelper[T, UserModel, IDTYPE]) getUserByUsername(ctx context.Context, username string, condition database.Condition[T]) (*UserModel, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}

	c := condition.New()
	for _, field := range u.loginValidationFields {
		c.Or(c.New().Set(field, database.ConditionOperationEqual, username))
	}

	users, err := u.CrudHelper.Get(ctx, append([]string{"id", "password"}, u.loginValidationFields...), c)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("username not found")
	}

	return &users[0], nil
}

func (u *userHelper[T, UserModel, IDTYPE]) UpdatePasswordFromToken(ctx context.Context, token, password string, condition database.Condition[T]) error {
	user, reason, err := JWTAuthValidate[UserModel]("Bearer "+token, u.secret)
	if err != nil {
		return err
	}

	if reason != "forget_password" {
		return fmt.Errorf("invalid token")
	}

	u.userValidator.SetPassword(user, password)
	condition = condition.New().Set("id", database.ConditionOperationEqual, (*user).GetID())

	return u.Update(ctx, user, []string{"password"}, condition)
}
