package main

import (
	"fmt"
	"time"
)

func main() {
	d := time.Now()
	_, err := GetAllPlatsRecordedSince(d)

	if err != nil {
		fmt.Println(err.Error())
	}
}