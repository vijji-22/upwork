package mysql

import (
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

type UserRoleMappingHelper database.MysqlCurlHelper[rbac.UserRoleMapping[int64]]

func NewUserRoleMappingHelper(db database.Connection) UserRoleMappingHelper {
	return database.NewBaseHelper[rbac.UserRoleMapping[int64]](db, "user_role_mapping", func(a *rbac.UserRoleMapping[int64]) map[string]interface{} {
		return map[string]interface{}{
			"id":      &a.ID,
			"user_id": &a.UserID,
			"role_id": &a.RoleID,

			"reference_id": &a.ReferenceID,
		}
	})
}
