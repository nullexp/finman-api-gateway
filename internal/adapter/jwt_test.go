package driven

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nullexp/finman-api-gateway/internal/port/model"
	"github.com/stretchr/testify/assert"
)

func TestNewTokenService(t *testing.T) {
	secret := "testsecret"
	expireAfter := time.Hour
	ts := NewTokenService(secret, expireAfter)

	assert.Equal(t, secret, ts.secret)
	assert.Equal(t, expireAfter, ts.expireAfter)
}

func TestTokenServiceCreateToken(t *testing.T) {
	secret := "testsecret"
	expireAfter := time.Hour
	ts := NewTokenService(secret, expireAfter)

	subject := model.Subject{UserId: uuid.New().String(), IsAdmin: true}
	token, err := ts.CreateToken(subject)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestTokenServiceGetToken(t *testing.T) {
	secret := "testsecret"
	expireAfter := time.Hour
	ts := NewTokenService(secret, expireAfter)

	subject := model.Subject{UserId: uuid.New().String(), IsAdmin: true}
	token, err := ts.CreateToken(subject)
	assert.NoError(t, err)

	claims, err := ts.GetToken(token)

	assert.NoError(t, err)

	sub, err := ts.GetSubject(claims.Subject)
	assert.NoError(t, err)
	assert.Equal(t, subject.UserId, sub.UserId)
	assert.Equal(t, subject.IsAdmin, sub.IsAdmin)
}

func TestTokenServiceCheckToken(t *testing.T) {
	secret := "testsecret"
	expireAfter := time.Hour
	ts := NewTokenService(secret, expireAfter)

	subject := model.Subject{UserId: uuid.New().String(), IsAdmin: true}
	token, err := ts.CreateToken(subject)
	assert.NoError(t, err)

	valid, err := ts.CheckToken(token)
	assert.NoError(t, err)
	assert.True(t, valid)
}

func TestTokenService_GetSubject(t *testing.T) {
	secret := "testsecret"
	expireAfter := time.Hour
	ts := NewTokenService(secret, expireAfter)

	subject := model.Subject{UserId: uuid.New().String(), IsAdmin: true}
	encodedSubject, err := json.Marshal(subject)
	assert.NoError(t, err)

	encodedSubjectString := base64.RawStdEncoding.EncodeToString(encodedSubject)
	decodedSubject, err := ts.GetSubject(encodedSubjectString)
	assert.NoError(t, err)
	assert.Equal(t, subject.UserId, decodedSubject.UserId)
	assert.Equal(t, subject.IsAdmin, decodedSubject.IsAdmin)
}

func TestTokenService_GetSubjectInvalidBase64(t *testing.T) {
	secret := "testsecret"
	expireAfter := time.Hour
	ts := NewTokenService(secret, expireAfter)

	invalidBase64 := "invalid_base64_string"
	_, err := ts.GetSubject(invalidBase64)
	assert.Error(t, err)
}

func TestTokenService_GetSubjectInvalidJSON(t *testing.T) {
	secret := "testsecret"
	expireAfter := time.Hour
	ts := NewTokenService(secret, expireAfter)

	invalidJSON := base64.RawStdEncoding.EncodeToString([]byte("invalid_json"))
	_, err := ts.GetSubject(invalidJSON)
	assert.Error(t, err)
}

func TestTokenService_CheckTokenInvalidSecret(t *testing.T) {
	secret := "testsecret"
	expireAfter := time.Hour
	ts := NewTokenService(secret, expireAfter)

	subject := model.Subject{UserId: uuid.New().String(), IsAdmin: true}
	token, err := ts.CreateToken(subject)
	assert.NoError(t, err)

	// Create a TokenService with a different secret
	invalidSecretTS := NewTokenService("invalidsecret", expireAfter)
	valid, err := invalidSecretTS.CheckToken(token)
	assert.Error(t, err)
	assert.False(t, valid)
}
