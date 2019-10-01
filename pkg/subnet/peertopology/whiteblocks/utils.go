package whiteblocks

import (
	"math/rand"
	"time"
)

// note: min/max are inclusive
func randBetween(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
