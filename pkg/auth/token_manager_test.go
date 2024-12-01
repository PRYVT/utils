package auth

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserUuidFromToken(t *testing.T) {
	// Set up the environment variable
	os.Setenv("SIGNING_SECRET", "testsecret")
	defer os.Unsetenv("SIGNING_SECRET")

	// Create a test UUID
	testUuid := uuid.New()
	tokenString, err := CreateToken(testUuid)
	assert.NoError(t, err)

	returnedUuid, err := GetUserUuidFromToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, testUuid, returnedUuid)
}
