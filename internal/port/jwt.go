package driven

import (
	"github.com/nullexp/finman-api-gateway/internal/port/model"
)

type Claims interface {
	Valid() error
}

type TokenService interface {
	CreateToken(sb model.Subject) (string, error)
	GetToken(tokenString string) (model.StandardClaims, error)
	CheckToken(tokenString string) (bool, error)
	GetSubject(subject string) (out model.Subject, err error)
}
