package mysql

import (
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

type RoleAccessMappingHelper database.MysqlCurlHelper[rbac.RoleAccessMapping[int64]]

func NewRoleAccessMappingHelper(db database.Connection) RoleAccessMappingHelper {
	return database.NewBaseHelper[rbac.RoleAccessMapping[int64]](db, "role_access_mapping", func(a *rbac.RoleAccessMapping[int64]) map[string]interface{} {
		return map[string]interface{}{
			"id":        &a.ID,
			"role_id":   &a.RoleID,
			"access_id": &a.AccessID,

			"project": &a.Project,
		}
	})
}
