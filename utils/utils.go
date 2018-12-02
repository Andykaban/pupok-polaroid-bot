package utils

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
)

func GetRandom(min int, max int) int {
	return rand.Intn(max - min) + min
}

func GetRandomFileName(tempDir string) string {
	randomPart := GetRandom(100, 10000000)
	randomFileName := fmt.Sprintf("image_%d.jpg", randomPart)
	return filepath.Join(tempDir, randomFileName)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
