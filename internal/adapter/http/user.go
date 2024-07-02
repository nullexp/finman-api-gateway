package http

import (
	"context"
	"net/http"
	"time"

	userv1 "github.com/nullexp/finman-api-gateway/internal/adapter/grpc/user/v1"
	"github.com/nullexp/finman-api-gateway/internal/port/model"
	httpapi "github.com/nullexp/finman-api-gateway/pkg/infrastructure/http/protocol"
	"github.com/nullexp/finman-api-gateway/pkg/infrastructure/http/protocol/model/openapi"
)

const PleaseReadTheErrorCode = "Please read the error message"
const UserBaseURL = "/users"

func NewUser(client userv1.UserServiceClient, parser model.SubjectParser) httpapi.Module {
	return UserHandler{client: client, parser: parser}
}

type UserHandler struct {
	client userv1.UserServiceClient
	parser model.SubjectParser
}

func (s UserHandler) GetRequestHandlers() []*httpapi.RequestDefinition {
	return []*httpapi.RequestDefinition{
		s.GetAllUsers(),
		s.PostUsers(),
		s.GetUserById(),
		s.UpdateUser(),
		s.DeleteUser(),
	}
}

func (s UserHandler) GetBaseURL() string {
	return UserBaseURL
}

const (
	UserManagement  = "User Management"
	UserDescription = "Use these APIs to access user resources"
)

func (s UserHandler) GetTag() openapi.Tag {
	return openapi.Tag{
		Name:        UserManagement,
		Description: UserDescription,
	}
}

func (s UserHandler) GetAllUsers() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:          "",
		Method:         http.MethodGet,
		FreeRoute:      false,
		AnyPermissions: []string{"ManageUsers"},
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &GetAllUsersResponse{},
			},
		},
		Handler: func(req httpapi.Request) {
			users, err := s.client.GetAllUsers(context.Background(), &userv1.GetAllUsersRequest{})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			req.Negotiate(http.StatusCreated, err, users)
		},
	}
}

func (s UserHandler) PostUsers() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:          "",
		Method:         http.MethodPost,
		FreeRoute:      false,
		Dto:            &CreateUserRequest{},
		AnyPermissions: []string{"ManageUsers"},
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &CreateUserResponse{},
			},
		},
		Handler: func(req httpapi.Request) {
			dto := req.MustGetDTO().(*CreateUserRequest)
			resp, err := s.client.CreateUser(context.Background(), &userv1.CreateUserRequest{
				Username: dto.Username,
				Password: dto.Password,
				RoleId:   dto.RoleId,
			})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			req.Negotiate(http.StatusCreated, err, CreateUserResponse{
				Id: resp.Id,
			})
		},
	}
}

func (s UserHandler) GetUserById() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:      "/:id",
		Method:     http.MethodGet,
		FreeRoute:  false,
		Dto:        &GetUserByIdRequest{},
		Parameters: simpleIdParamDef,
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &GetUserByIdResponse{},
			},
		},
		Handler: func(req httpapi.Request) {
			id := req.MustGet(idDef.GetName()).(string)
			resp, err := s.client.GetUserById(context.Background(), &userv1.GetUserByIdRequest{
				Id: id,
			})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			req.Negotiate(http.StatusOK, nil, GetUserByIdResponse{
				User: UserReadable{
					Id:        resp.User.Id,
					Username:  resp.User.Username,
					RoleId:    resp.User.RoleId,
					IsAdmin:   resp.User.IsAdmin,
					CreatedAt: MustParseTime(resp.User.CreatedAt),
					UpdatedAt: MustParseTime(resp.User.UpdatedAt),
				},
			})
		},
	}
}

func (s UserHandler) UpdateUser() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:          "/:id",
		Method:         http.MethodPut,
		FreeRoute:      false,
		Dto:            &UpdateUserRequest{},
		AnyPermissions: []string{"ManageUsers"},
		Parameters:     simpleIdParamDef,
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &GetUserByIdResponse{},
			},
		},
		Handler: func(req httpapi.Request) {
			id := req.MustGet(idDef.GetName()).(string)
			dto := req.MustGetDTO().(*UpdateUserRequest)
			_, err := s.client.UpdateUser(context.Background(), &userv1.UpdateUserRequest{
				Id:       id,
				Password: dto.Password,
				RoleId:   dto.RoleId,
			})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			req.Negotiate(http.StatusOK, nil, nil)
		},
	}
}

func (s UserHandler) DeleteUser() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:          "/:id",
		Method:         http.MethodDelete,
		FreeRoute:      false,
		Dto:            &DeleteUserRequest{},
		AnyPermissions: []string{"ManageUsers"},
		Parameters:     simpleIdParamDef,
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         nil,
			},
		},
		Handler: func(req httpapi.Request) {
			id := req.MustGet(idDef.GetName()).(string)
			_, err := s.client.DeleteUser(context.Background(), &userv1.DeleteUserRequest{
				Id: id,
			})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			req.ReturnStatus(http.StatusOK, nil)
		},
	}
}

type UserReadable struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	RoleId    string    `json:"roleId"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,gte=1"`
	Password string `json:"password" validate:"required,gte=1"`
	RoleId   string `json:"roleId" validate:"required,uuid"`
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
	RoleId   string `json:"roleId" validate:"required,uuid"`
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
