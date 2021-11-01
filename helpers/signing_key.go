package helpers

import (
	"log"
	"os"
	"sync"
)

var onlyOnce sync.Once

//        /!\ NOT SAFE /!\        //
const DEFAULT_PRIVATE_KEY string = "nM8g5JYoVZJtHnfbySKDq8vBqErrmz0AqleeWXKcoMRXsMoYpHHg8blxRilXHkuj"

// GetSigningKey gets private key from environment variable
func GetSigningKey() string {
	privKey := os.Getenv("PRIVATE_KEY")

	if privKey == "" {
		onlyOnce.Do(func() {
			log.Println(`/!\/!\/!\ Your API is using a hardcoded private key /!\/!\/!\`)
		})

		return DEFAULT_PRIVATE_KEY
	}

	return privKey
}
