package util

import (
	"crypto/sha512"
	"fmt"
)

func GetSHA512Hash(input string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(input)))
}
