package helpers

import (
	"math/rand"
	"strings"
	"time"
)

const ALPHANUM string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateID(nb int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	var sb strings.Builder
	sb.Grow(nb)

	for ; nb > 0; nb-- {
		sb.WriteByte(ALPHANUM[rand.Intn(len(ALPHANUM)-1)])
	}

	return sb.String()
}
