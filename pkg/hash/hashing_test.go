package hash

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"
	hashedPassword := HashPassword(password)

	if hashedPassword != "94aefb8be78b2b7c344d11d1ba8a79ef087eceb19150881f69460b8772753263" {
		t.Errorf("Expected hashed password to be non-empty")
	}

}
