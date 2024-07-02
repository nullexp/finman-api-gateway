package model

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type StandardClaims struct {
	Audience  []string `json:"aud,omitempty"`
	ExpiresAt int64    `json:"exp,omitempty"`
	Identity  string   `json:"jti,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	Issuer    string   `json:"iss,omitempty"`
	NotBefore int64    `json:"nbf,omitempty"`
	Subject   string   `json:"sub,omitempty"`
}

func (c StandardClaims) Valid() error {
	return nil
}

func (c StandardClaims) GetExpireTime() int64 {
	return c.ExpiresAt
}

func (c StandardClaims) GetSubject() string {
	return c.Subject
}

func (c StandardClaims) GetIssuer() string {
	return c.Issuer
}

func (c StandardClaims) GetAudience() []string {
	return c.Audience
}

func (c StandardClaims) GetIssuedAt() int64 {
	return c.IssuedAt
}

func (c StandardClaims) GetIdentity() string {
	return c.Identity
}

func (c StandardClaims) IsExpired() bool {
	return time.Now().Unix() > c.ExpiresAt
}

type SubjectParser interface {
	MustParseSubject(string) Subject
}

type Subject struct {
	UserId  string `json:"userId"`
	IsAdmin bool   `json:"isAdmin"`
}

type testSubjectParser struct {
	subject Subject
}

func NewTestSubjectParser(s Subject) testSubjectParser {
	return testSubjectParser{subject: s}
}

func (ts testSubjectParser) MustParseSubject(string) Subject {
	return ts.subject
}

func ToSubject(subject string) (out Subject, err error) {
	data, err := base64.RawStdEncoding.DecodeString(subject)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(data, &out)
	return
}
