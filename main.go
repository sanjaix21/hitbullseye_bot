package main

import (
	"hitbullseye_am/internal"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	hitbullseyeId := os.Getenv("HBE_ID")
	hitbullseyePassword := os.Getenv("HBE_PASS")

	var err error
	handler := internal.NewHandler(hitbullseyeId, hitbullseyePassword)

	err = handler.Login()
	if err != nil {
		panic(err)
	}

	err = handler.NavigateToTest()
	if err != nil {
		panic(err)
	}

	err = handler.GiveTest()
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Hour)
}
