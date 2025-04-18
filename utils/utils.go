package utils

import (
	"fmt"
	"math/rand"
)

func GenerateNoRek() string {
	randomNumber := rand.Int63n(10000000000)
	return fmt.Sprintf("%010d", randomNumber)
}
