package model

import (
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac/mysql"
)

type Template struct {
	database.TableID[int64]
	Name        string                              `json:"name,omitempty"`
	DataMapping database.DbMap[string, interface{}] `json:"datamapping,omitempty"`
	FilePath    string                              `json:"filepath,omitempty"`
}

func GetTemplateHelper(db database.Connection) database.CrudHelper[database.MysqlCondition, Template, int64] {
	h := database.NewBaseHelper(db, "template", func(t *Template) map[string]interface{} {
		return map[string]interface{}{
			"id":          &t.ID,
			"name":        &t.Name,
			"datamapping": &t.DataMapping,
			"filepath":    &t.FilePath,
		}
	})

	rbacHelper := rbac.NewRbacHelper(mysql.MysqlRbacHelper(db))
	rbacCrudHelper := rbac.NewCrudHelper(rbacHelper, h, UserIDFromCTX)

	return rbacCrudHelper
}
