// Package gen implements different random values
// generators.
package gen

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

func String(length int) (string, error) {
	s := new(strings.Builder)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", fmt.Errorf("rand.Int: %w", err)
		}

		if err := s.WriteByte(letters[n.Int64()]); err != nil {
			return "", fmt.Errorf("s.WriteByte: %w", err)
		}
	}

	return s.String(), nil
}
