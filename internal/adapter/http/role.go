package http

import (
	"context"
	"net/http"
	"time"

	userv1 "github.com/nullexp/finman-api-gateway/internal/adapter/grpc/user/v1"
	httpapi "github.com/nullexp/finman-api-gateway/pkg/infrastructure/http/protocol"
	"github.com/nullexp/finman-api-gateway/pkg/infrastructure/http/protocol/model/openapi"
)

const RoleBaseURL = "/roles"

func NewRole(client userv1.RoleServiceClient) httpapi.Module {
	return RoleHandler{client: client}
}

type RoleHandler struct {
	client userv1.RoleServiceClient
}

func (s RoleHandler) GetRequestHandlers() []*httpapi.RequestDefinition {
	return []*httpapi.RequestDefinition{
		s.GetAllRoles(),
	}
}

func (s RoleHandler) GetBaseURL() string {
	return RoleBaseURL
}

const (
	RoleManagement  = "Role Management"
	RoleDescription = "Use these apis to access role resources"
)

func (s RoleHandler) GetTag() openapi.Tag {
	return openapi.Tag{
		Name:        RoleManagement,
		Description: RoleDescription,
	}
}

func (s RoleHandler) GetAllRoles() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:          "",
		Method:         http.MethodGet,
		FreeRoute:      false,
		AnyPermissions: []string{"ManageRoles"},
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &GetAllRolesResponse{},
			},
		},
		Handler: func(req httpapi.Request) {
			Roles, err := s.client.GetAllRoles(context.Background(), &userv1.GetAllRolesRequest{})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			req.Negotiate(http.StatusCreated, err, Roles)
		},
	}
}

func (s RoleHandler) PostRoles() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:          "",
		Description:    "Please note that permissions can be these: ManageRoles ManageUsers ManageTransactions",
		Method:         http.MethodPost,
		FreeRoute:      false,
		Dto:            &CreateRoleRequest{},
		AnyPermissions: []string{"ManageRoles"},
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &CreateRoleResponse{},
			},
		},
		Handler: func(req httpapi.Request) {
			dto := req.MustGetDTO().(*CreateRoleRequest)
			resp, err := s.client.CreateRole(context.Background(), &userv1.CreateRoleRequest{
				Name:        dto.Name,
				Permissions: dto.Permissions,
			})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			req.Negotiate(http.StatusCreated, err, CreateRoleResponse{
				Id: resp.Id,
			})
		},
	}
}

type Role struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required"`
	Permissions []string `json:"permissions" validate:"required"`
}

func (dto CreateRoleRequest) Validate(ctx context.Context) error {
	return nil
}

type CreateRoleResponse struct {
	Id string `json:"id"`
}

type GetRoleByIdRequest struct {
	Id string `json:"id" validate:"required,uuid"`
}

func (dto GetRoleByIdRequest) Validate(ctx context.Context) error {
	return nil
}

type GetRoleByIdResponse struct {
	Role Role `json:"role"`
}

type GetAllRolesResponse struct {
	Roles []Role `json:"roles"`
}

type UpdateRoleRequest struct {
	Id          string   `json:"id" validate:"required,uuid"`
	Name        string   `json:"name" validate:"required"`
	Permissions []string `json:"permissions" validate:"required"`
}

func (dto UpdateRoleRequest) Validate(ctx context.Context) error {
	return nil
}

type DeleteRoleRequest struct {
	Id string `json:"id" validate:"required,uuid"`
}

func (dto DeleteRoleRequest) Validate(ctx context.Context) error {
	return nil
}

type IsUserPermittedToPermissionRequest struct {
	UserId     string `json:"userId" validate:"required,uuid"`
	Permission string `json:"permission" validate:"required"`
}

func (dto IsUserPermittedToPermissionRequest) Validate(ctx context.Context) error {
	return nil
}

type IsUserPermittedToPermissionResponse struct {
	IsPermitted bool `json:"isPermitted"`
}
