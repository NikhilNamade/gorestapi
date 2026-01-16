package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)
const secretKey = "Thisiseventmanagement"

func GenerateToken(userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
	})
	return token.SignedString([]byte(secretKey))
}

func AuthenticateUser(token string) (int64,error){
	parsedToken,err := jwt.Parse(token,func(t *jwt.Token) (any, error) {
		_,ok:=t.Method.(*jwt.SigningMethodHMAC)

		if !ok{
			return  nil, errors.New("Unexpected signing method")
		}
		return []byte(secretKey),nil
	})


	if err != nil{
		return 0,err
	}
	if !parsedToken.Valid{
		return 0,errors.New("Invalid token!")
	}

	claims,ok:=parsedToken.Claims.(jwt.MapClaims)

	if !ok{
		return 0,errors.New("Invalid token!")
	}

	userId := claims["userId"].(float64)

	return int64(userId),nil
}
