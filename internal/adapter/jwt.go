package driven

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	driven "github.com/nullexp/finman-api-gateway/internal/port"
	"github.com/nullexp/finman-api-gateway/internal/port/model"
)

// TokenService is a struct that manages JWT tokens.
type TokenService struct {
	secret      string
	expireAfter time.Duration
}

// NewTokenService creates a new TokenService with the provided secret.
func NewTokenService(secret string, expireAfter time.Duration) *TokenService {
	return &TokenService{secret: secret, expireAfter: expireAfter}
}

// CreateToken generates a JWT token for the given subject.
func (ts TokenService) CreateToken(sb model.Subject) (string, error) {
	// Marshal the subject to JSON.
	data, err := json.Marshal(sb)
	if err != nil {
		log.Printf("Error marshaling subject: %v", err)
		return "", err
	}

	// Encode the JSON data to a base64 string.
	enc := base64.RawStdEncoding.EncodeToString(data)
	log.Printf("Encoded subject to base64: %s", enc)

	// Create the token with the encoded subject.
	return ts.createTokenWithText(enc, ts.expireAfter)
}

// CreateTokenWithText generates a JWT token with the provided text and expireTime.
func (ts TokenService) createTokenWithText(sb string, expireAfter time.Duration) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("HS256"))
	t.Claims = model.StandardClaims{Subject: sb, ExpiresAt: time.Now().Add(expireAfter).Unix(), Identity: uuid.NewString()}

	// Sign the token with the secret.
	tokenString, err := t.SignedString([]byte(ts.secret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}

	return tokenString, nil
}

// GetToken parses the given token string and returns the claims.
func (ts TokenService) GetToken(tokenString string) (model.StandardClaims, error) {
	sc := model.StandardClaims{}

	// Parse the token.
	rawToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(ts.secret), nil
	})
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return sc, err
	}

	// Split the token into parts.
	parts := strings.Split(rawToken.Raw, ".")
	if len(parts) != 3 {
		err = errors.New("invalid token format")
		log.Printf("Error splitting token parts: %v", err)
		return sc, err
	}

	// Decode the base64 part of the token.
	data, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("Error decoding base64 token part: %v", err)
		return sc, err
	}

	// Unmarshal the JSON data into the standard claims.
	err = json.Unmarshal(data, &sc)
	if err != nil {
		err = errors.New("unknown subject")
		log.Printf("Error unmarshaling JSON data: %v", err)
		return sc, err
	}

	log.Printf("Parsed token claims: %+v", sc)
	return sc, nil
}

// CheckToken validates the given token string.
func (ts TokenService) CheckToken(tokenString string) (bool, error) {
	// Parse the token.
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New("unexpected signing method")
		}
		return []byte(ts.secret), nil
	})
	// Check if there was an error parsing the token.
	if err != nil {
		log.Printf("Error checking token: %v", err)
		return false, err
	}

	log.Printf("Token is valid")
	return true, nil
}

func (ts TokenService) GetSubject(subject string) (out model.Subject, err error) {
	data, err := base64.RawStdEncoding.DecodeString(subject)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(data, &out)
	return
}

// NewValidJwtClaim creates a new valid JWT claim with the given expiration time.
func NewValidJwtClaim(expireTime time.Duration) driven.JwtClaim {
	return model.StandardClaims{ExpiresAt: time.Now().Add(expireTime).Unix()}
}
