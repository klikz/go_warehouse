package apiserver

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"warehouse/internal/app/models"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
func isPasswordValid(p string) bool {
	return len(p) > 6
}

func ComparePassword(password, encrypt string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encrypt), []byte(password)) == nil
}

func GetToken(u *models.User) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role":  u.Role,
		"email": u.Email,
		"nbf":   time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	tokenString, err := token.SignedString(Secret_key)
	if err != nil {
		return err
	}
	u.Token = tokenString
	return nil
}

func ParseToken(tokenString string) (models.ParsedToken, error) {
	parsedToken := models.ParsedToken{}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return Secret_key, nil
	})
	if err != nil {
		return parsedToken, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

	} else {
		return parsedToken, errors.New("wrong token")
	}
	parsedToken.Role = fmt.Sprint(claims["role"])
	parsedToken.Email = fmt.Sprint(claims["email"])

	return parsedToken, nil
}