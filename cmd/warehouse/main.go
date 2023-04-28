package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"warehouse/internal/app/apiserver"

	"github.com/BurntSushi/toml"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jasonlvhit/gocron"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func checkLicenseTime() {
	b, err := os.ReadFile("key")
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	token2, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("MusayevSDRuffrider123SomeSecretKeys"), nil
	})
	regTime1 := ""
	regTime2 := ""
	if claims, ok := token2.Claims.(jwt.MapClaims); ok && token2.Valid {
		// fmt.Println("date: ", claims["date"])
		regTime1 = fmt.Sprint(claims["date1"])
		regTime2 = fmt.Sprint(claims["date2"])
	} else {
		log.Fatal(err)
		fmt.Println(err)
	}

	t1, _ := time.Parse(time.RFC3339, regTime1)
	t2, _ := time.Parse(time.RFC3339, regTime2)

	// fmt.Println("t: ", t.Sub(time.Now()).Hours())

	elapsed1 := t1.Sub(time.Now()).Hours()
	elapsed2 := t2.Sub(time.Now()).Hours()

	// fmt.Println("elapsed1: ", elapsed1)
	fmt.Println(elapsed2)

	if elapsed1 < -720 {
		log.Fatal("Time LIMIT exceed")
	}

	if elapsed2 < 0 {
		log.Fatal("Time LIMIT exceed")
	}
}

func executeCronJob() {
	gocron.Every(5).Second().Do(checkLicenseTime)
	<-gocron.Start()
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	checkLicenseTime()
	go executeCronJob()

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}

}
