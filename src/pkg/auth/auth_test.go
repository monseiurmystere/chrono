package auth

import (
	"testing"
)

func TestClearSecret(t *testing.T) {
	secret := []byte("super-secret-password-1234")
	ClearSecret(secret)

	for i, b := range secret {
		if b != 0 {
			t.Fatalf("expected byte at index %d to be 0, got %d", i, b)
		}
	}
}
