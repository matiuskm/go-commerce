package helpers

import (
	"fmt"
	"math/rand"
	"time"
)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateOrderNumber() string {
	datePart := time.Now().Format("20060102")
	randomPart := generateRandomString(6)
	return fmt.Sprintf("ORD-%s-%s", datePart, randomPart)
}