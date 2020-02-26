package shitpost

import (
	"fmt"
	"math/rand"
	"time"
)

func brexitCountdown() string {
	timeBrexit, _ := time.Parse(time.RFC3339, "2020-01-31T23:00:00Z") // it literally never errors
	d := time.Since(timeBrexit)
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	lines := []string{
		fmt.Sprintf("%d days, %d hours, %d minutes and %d seconds since the British Empire has be freed from the oppression of the EU", days, hours, minutes, seconds),
		fmt.Sprintf("%d days, %d hours, %d minutes and %d seconds since Great Britain collapsed", days, hours, minutes, seconds),
	}
	rand.Seed(time.Now().Unix())
	return lines[rand.Intn(len(lines))]
}
