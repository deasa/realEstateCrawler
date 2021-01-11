package main

import (
	"fmt"
	"time"
)

func main() {
	d := time.Now()
	plats, err := GetAllPlatsRecordedSince(d)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}

	for _, plat := range plats {
		fmt.Println(plat)
	}
}