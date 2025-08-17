package main

import (
	"fmt"
	"hitbullseye_bot/internal"
	"os"
	"time"

	"github.com/subosito/gotenv"
)

func main() {
	_ = gotenv.Load()
	var hitbullseyeId, hitbullseyePassword string
	hitbullseyeId = os.Getenv("HBE_ID")
	hitbullseyePassword = os.Getenv("HBE_PASS")

	if hitbullseyeId == "" || hitbullseyePassword == "" {
		fmt.Println("⚠️  Please open your .env file and add:\nHBE_ID=your_id\nHBE_PASS=your_pass")
		return
	}
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
