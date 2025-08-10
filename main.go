package main

import (
	"hitbullseye_am/internal"
	"time"
)

func main() {
	hitbullseyeId := "100026942"
	hitbullseyePassword := "73696062"

	handler := internal.NewHandler(hitbullseyeId, hitbullseyePassword)
	err := handler.Login()
	if err != nil {
		panic(err)
	}

	err = handler.NavigateToTest()
	if err != nil {
		panic(err)
	}

	// time.Sleep(time.Second * 5)
	// fmt.Println("Done sleeping")
	err = handler.GiveTest()
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Hour)
}
