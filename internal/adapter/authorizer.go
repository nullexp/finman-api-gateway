package adapter

import (
	"context"

	userv1 "github.com/nullexp/finman-api-gateway/internal/adapter/grpc/user/v1"
	"github.com/nullexp/finman-api-gateway/internal/port/model"
	"github.com/nullexp/finman-api-gateway/pkg/infrastructure/http/protocol"
)

func NewAuthorizer(client userv1.RoleServiceClient, parser model.SubjectParser) protocol.Authorizer {
	return func(identity, permission string) (bool, error) {

		sub := parser.MustParseSubject(identity)
		if sub.IsAdmin {
			return true, nil
		}
		rs, err := client.IsUserPermittedToPermission(context.Background(), &userv1.IsUserPermittedToPermissionRequest{
			UserId:     sub.UserId,
			Permission: permission,
		})

		if err != nil {
			return false, err
		}
		return rs.IsPermitted, nil
	}
}
