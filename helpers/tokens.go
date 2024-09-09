package helpers

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "1234534sdgfdsbdvccxvfsdgf"

func GenerateToken(email string, userId int64, isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   email,
		"userId":  userId,
		"isAdmin": isAdmin,
		"exp":     time.Now().Add(time.Hour * 8766).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenString string) (int64, bool, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, OK := token.Method.(*jwt.SigningMethodHMAC)
		if !OK {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, false, errors.New("could not parse token")
	}
	IsValidToken := parsedToken.Valid

	if !IsValidToken {
		return 0, false, errors.New("invalid token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false, errors.New("could not parse claims")
	}

	userId := int64(claims["userId"].(float64))
	isAdmin := claims["isAdmin"].(bool)

	return userId, isAdmin, nil
}
