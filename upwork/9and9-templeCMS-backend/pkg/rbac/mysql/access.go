package mysql

import (
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

type AccessHelper database.MysqlCurlHelper[rbac.Access[int64]]

func NewAccessHelper(db database.Connection) AccessHelper {
	return database.NewBaseHelper[rbac.Access[int64]](db, "access", func(a *rbac.Access[int64]) map[string]interface{} {
		return map[string]interface{}{
			"id":            &a.ID,
			"name":          &a.Name,
			"reference_key": &a.ReferenceKey,
		}
	})
}
