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

	fmt.Printf("current time: %s\n", formatTime(time.Now()))
	fmt.Printf("exact time: %s\n", formatTime(ntpTime))
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05 +0000 UTC")
}
