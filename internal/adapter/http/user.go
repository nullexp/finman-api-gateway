package http

import (
	"context"
	"log"
	"net/http"
	"time"

	userv1 "github.com/nullexp/finman-api-gateway/internal/adapter/grpc/user/v1"
	"github.com/nullexp/finman-api-gateway/internal/port/model"
	httpapi "github.com/nullexp/finman-api-gateway/pkg/infrastructure/http/protocol"
	"github.com/nullexp/finman-api-gateway/pkg/infrastructure/http/protocol/model/openapi"
)

const UserBaseURL = "/users"

func NewUser(client userv1.UserServiceClient, parser model.SubjectParser) httpapi.Module {
	return User{client: client, parser: parser}
}

type User struct {
	client userv1.UserServiceClient
	parser model.SubjectParser
}

func (s User) GetRequestHandlers() []*httpapi.RequestDefinition {
	return []*httpapi.RequestDefinition{
		s.GetAllUsers(),
	}
}

func (s User) GetBaseURL() string {
	return UserBaseURL
}

const (
	UserManagement  = "User Management"
	UserDescription = "Use these apis to access user resources"
)

func (s User) GetTag() openapi.Tag {
	return openapi.Tag{
		Name:        UserManagement,
		Description: UserDescription,
	}
}

func (s User) GetAllUsers() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:     "",
		Method:    http.MethodPost,
		FreeRoute: false,
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &GetAllUsersResponse{},
			},
		},
		Handler: func(req httpapi.Request) {
			caller := req.MustGetCaller()
			sub := s.parser.MustParseSubject(caller.GetSubject())
			log.Println(sub)
			users, err := s.client.GetAllUsers(context.Background(), &userv1.GetAllUsersRequest{})
			req.Negotiate(http.StatusCreated, err, users)
		},
	}
}

type UserReadable struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	RoleId    string    `json:"role_id"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,gte=1"`
	Password string `json:"password" validate:"required,gte=1"`
	RoleId   string `json:"role_id" validate:"required,uuid"`
}

func (dto CreateUserRequest) Validate(ctx context.Context) error {
	return nil
}

type CreateUserResponse struct {
	Id string `json:"id"`
}
type GetUserByIdRequest struct {
	Id string `json:"id" validate:"required,uuid"`
}

func (dto GetUserByIdRequest) Validate(ctx context.Context) error {
	return nil
}

type GetUserByIdResponse struct {
	User UserReadable `json:"user"`
}
type GetAllUsersResponse struct {
	Users []UserReadable `json:"users"`
}
type UpdateUserRequest struct {
	Id       string `json:"id" validate:"required,uuid"`
	Password string `json:"password" validate:"required,gte=1"`
	RoleId   string `json:"role_id" validate:"required,uuid"`
}

func (dto UpdateUserRequest) Validate(ctx context.Context) error {
	return nil
}

type DeleteUserRequest struct {
	Id string `json:"id" validate:"required,uuid"`
}

func (dto DeleteUserRequest) Validate(ctx context.Context) error {
	return nil
}

type GetUserByUsernameAndPasswordRequest struct {
	Username string `json:"username" validate:"required,gte=1"`
	Password string `json:"password" validate:"required,gte=1"`
}

func (dto GetUserByUsernameAndPasswordRequest) Validate(ctx context.Context) error {
	return nil
}

type GetUsersWithPaginationRequest struct {
	Limit  int `json:"limit" validate:"gte=0"`
	Offset int `json:"offset" validate:"gte=0"`
}

func (dto GetUsersWithPaginationRequest) Validate(ctx context.Context) error {
	return nil
}

type GetUserByUsernameAndPasswordResponse struct {
	User UserReadable `json:"user"`
}

type GetUsersWithPaginationResponse struct {
	Users []UserReadable `json:"users"`
}
