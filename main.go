package main

import (
	"hitbullseye_am/internal"
	"time"
)

func main() {
	hitbullseyeId := "100027311"
	hitbullseyePassword := "98322853"

	var err error
	handler := internal.NewHandler(hitbullseyeId, hitbullseyePassword)

	// err = handler.Login()
	// if err != nil {
	// 	panic(err)
	// }

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
