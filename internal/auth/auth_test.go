package auth
import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "securePassword123"
	
	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
	
	// Verify the hash works for comparison
	err = CheckPasswordHash(hash, password)
	assert.NoError(t, err)
	
	// Verify wrong password fails
	err = CheckPasswordHash(hash, "wrongPassword")
	assert.Error(t, err)
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	duration := 1 * time.Hour
	
	// Create a token
	token, err := MakeJWT(userID, tokenSecret, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Validate the token
	parsedUserID, err := ValidateJWT(token, tokenSecret)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

func TestExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	
	// Create a token that expires immediately
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
		Subject:   userID.String(),
	})
	
	signedToken, err := token.SignedString([]byte(tokenSecret))
	assert.NoError(t, err)
	
	// Attempt to validate the expired token
	_, err = ValidateJWT(signedToken, tokenSecret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	correctSecret := "correct-secret"
	wrongSecret := "wrong-secret"
	duration := 1 * time.Hour
	
	// Create a token with the correct secret
	token, err := MakeJWT(userID, correctSecret, duration)
	assert.NoError(t, err)
	
	// Validate with wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature is invalid")
}

func TestInvalidToken(t *testing.T) {
	// Test with a completely invalid token
	_, err := ValidateJWT("not-a-valid-token", "any-secret")
	assert.Error(t, err)
}

func TestInvalidUUID(t *testing.T) {
	tokenSecret := "test-secret-key"
	
	// Create a token with an invalid UUID as subject
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		Subject:   "not-a-valid-uuid",
	})
	
	signedToken, err := token.SignedString([]byte(tokenSecret))
	assert.NoError(t, err)
	
	// Validate should fail due to invalid UUID
	_, err = ValidateJWT(signedToken, tokenSecret)
	assert.Error(t, err)
}