package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"date1": time.Now().Local().Add(time.Hour * time.Duration(0)),
		"date2": time.Now().Local().Add(time.Hour * time.Duration(720)),
		"nbf":   time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenStringTemp, err := token.SignedString([]byte("MusayevSDRuffrider123SomeSecretKeys"))
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("key", []byte(tokenStringTemp), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
