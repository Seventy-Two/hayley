package shitpost

import (
	"fmt"
	"math/rand"
	"time"
)

func brexitCountdown() string {
	timeBrexit, _ := time.Parse(time.RFC3339, "2019-06-30T23:00:00Z") // it literally never errors
	d := time.Until(timeBrexit)
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	lines := []string{
		fmt.Sprintf("%d days and %d hours until the British Empire is free from the oppression of the EU", days, hours),
		fmt.Sprintf("%d days and %d hours until Great Britain collapses", days, hours),
	}
	rand.Seed(time.Now().Unix())
	return lines[rand.Intn(len(lines))]
}
