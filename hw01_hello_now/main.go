package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const TimeSourceURL = "0.beevik-ntp.pool.ntp.org"

func main() {
	ntpTime, err := ntp.Time(TimeSourceURL)

	if err != nil {
		log.Fatalf("ntp error: %s", err.Error())
	}

	fmt.Printf("current time: %s\n", time.Now().Round(0))
	fmt.Printf("exact time: %s\n", ntpTime.Round(0))
}
