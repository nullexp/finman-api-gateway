package http

import (
	"context"
	"net/http"

	authv1 "github.com/nullexp/finman-api-gateway/internal/adapter/grpc/auth/v1"
	httpapi "github.com/nullexp/finman-api-gateway/pkg/infrastructure/http/protocol"
	"github.com/nullexp/finman-api-gateway/pkg/infrastructure/http/protocol/model/openapi"
)

const SessionBaseURL = "/sessions"

func NewSession(client authv1.AuthServiceClient) httpapi.Module {
	return Session{client: client}
}

type Session struct {
	client authv1.AuthServiceClient
}

func (s Session) GetRequestHandlers() []*httpapi.RequestDefinition {
	return []*httpapi.RequestDefinition{
		s.PostSession(),
	}
}

func (s Session) GetBaseURL() string {
	return SessionBaseURL
}

const (
	SessionManagement  = "Session Management"
	SessionDescription = "Authenticate through these apis"
)

func (s Session) GetTag() openapi.Tag {
	return openapi.Tag{
		Name:        SessionManagement,
		Description: SessionDescription,
	}
}

const RouteRefresh = "/refresh"

func (s Session) PostSession() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:     "",
		Dto:       &CreateTokenRequest{},
		FreeRoute: true,
		Method:    http.MethodPost,
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusCreated,
				Description: "If auth info is valid",
				Dto:         &CreateTokenResponse{},
			},
			{
				Status:      http.StatusBadRequest,
				Description: "If auth info is not valid",
			},
		},
		Handler: func(req httpapi.Request) {
			dto := req.MustGetDTO().(*CreateTokenRequest)
			token, err := s.client.Login(context.Background(), &authv1.LoginRequest{Username: dto.Username, Password: dto.Password})
			req.Negotiate(http.StatusCreated, err, token)
		},
	}
}

type CreateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (CreateTokenRequest) Validate(context.Context) error { return nil }

type CreateTokenResponse struct {
	Token string `json:"token"`
}
