package utils

import (
	"math/rand"
	"time"
)

func GetRandom(min int, max int) int {
	return rand.Intn(max - min) + min
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
