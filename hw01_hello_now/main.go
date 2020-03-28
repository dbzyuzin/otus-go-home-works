package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"log"
	"time"
)

const TimeSourceUrl = "0.beevik-ntp.pool.ntp.org"

func main() {
	ntpTime, err := ntp.Time(TimeSourceUrl)

	if err != nil {
		log.Fatalf("ntp error: %s", err.Error())
	}

	fmt.Printf("current time: %s\n", time.Now().Round(0))
	fmt.Printf("exact time: %s\n", ntpTime.Round(0))
}
