package whiteblocks

import "math/rand"

func randBetween(min, max int) int {
	return rand.Intn(max-min) + min
}
