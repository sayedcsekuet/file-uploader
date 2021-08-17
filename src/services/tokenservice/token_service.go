package tokenservice

import (
	"file-uploader/src/errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenService interface {
	Generate(fileId, key string, exp int64) (string, error)
	Verify(tokenStr, fileId, key string) error
}
type tokenService struct {
}

const maxDurationHour = 87600 //30 years

func NewTokenService() TokenService {
	return &tokenService{}
}

func (jm *tokenService) Generate(fileId, key string, exp int64) (string, error) {
	if exp == 0 {
		exp = time.Now().Add(time.Hour * maxDurationHour).Unix()
	}
	claims := jwt.MapClaims{"authorized": true, "file_id": fileId, "exp": exp}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(fmt.Sprintf("%s%s", key, fileId)))
	if err != nil {
		return "", err
	}
	return token, nil
}
func (jm *tokenService) Verify(tokenStr, fileId, key string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(fmt.Sprintf("%s%s", key, fileId)), nil
	})
	if err != nil {
		return errors.NewKnown(400, err.Error())
	}
	_, ok := token.Claims.(jwt.Claims)
	if !ok && !token.Valid {
		return errors.NewKnown(403, "Token is not valid!")
	}
	return nil
}
