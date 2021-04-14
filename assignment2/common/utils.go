package common

import (
	"log"
	"time"
)

// Usage: defer timeTrack(time.Now(), "roulette") before the block you want to time
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
