package services

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/matthewhartstonge/argon2"
	"github.com/rs/zerolog"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/models"
	"github.com/uptrace/bun"
)

type UserService struct {
	db     *bun.DB
	logger *zerolog.Logger
	conf   *core.Config
}

func NewUserService(db *bun.DB, logger *zerolog.Logger, conf *core.Config) *UserService {
	return &UserService{db, logger, conf}
}

func (u *UserService) UserExists(ctx context.Context, email string) bool {
	user := new(models.User)

	count, _ := u.db.NewSelect().Model(user).Where("email = ?", email).Count(ctx)

	return count != 0
}

func (u *UserService) CreateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": user})
	tokenString, err := token.SignedString([]byte(u.conf.SecretKey))

	return tokenString, err
}

func (u *UserService) HashPassword(password string) (string, error) {
	argon := argon2.RecommendedDefaults()
	argon.Mode = argon2.ModeArgon2i

	hash, err := argon.Hash([]byte(password), []byte(u.conf.SecretKey))
	if err != nil {
		return "", err
	}

	return string(hash.Encode()), nil
}

func (u *UserService) CheckPassword(user *models.User, password string) (bool, error) {
	return argon2.VerifyEncoded([]byte(password), []byte(user.PasswordHash))
}
