package mysql

import (
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

type RoleHelper database.MysqlCurlHelper[rbac.Role[int64]]

func NewRoleHelper(db database.Connection) RoleHelper {
	return database.NewBaseHelper[rbac.Role[int64]](db, "role", func(r *rbac.Role[int64]) map[string]interface{} {
		return map[string]interface{}{
			"id":   &r.ID,
			"name": &r.Name,
		}
	})
}
