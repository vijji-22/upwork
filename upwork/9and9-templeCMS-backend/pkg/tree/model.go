package tree

import (
	"context"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

/*
CREATE TABLE IF NOT EXISTS ` + tableName + ` (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			parent_id int(11) unsigned,
			group_id int(11) unsigned NOT NULL,
			node_type_id int(11) unsigned NOT NULL,
			value varchar(255) NOT NULL,
			created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS tree_fields (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			parent_id int(11) unsigned,
			config_id int(11) unsigned NOT NULL,
			name varchar(255) NOT NULL,
			type varchar(255) NOT NULL,
			required boolean NOT NULL,
			created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY (name, parent_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
*/

type ValueType string

const (
	String  ValueType = "string"
	Number  ValueType = "number"
	Boolean ValueType = "boolean"
	Array   ValueType = "array"
	Object  ValueType = "object"
)

type dbNodemeta struct {
	database.TableID[int64]
	ParentID *int64 `json:"parent_id,omitempty"`
	ConfigID int64  `json:"config_id"`

	Name string    `json:"name,omitempty"`
	Type ValueType `json:"type,omitempty"`

	Required bool `json:"required"`
}

type dbNodemetaHelper struct {
	*database.BaseHelper[dbNodemeta]
	configID int64
}

func NewDbNodemetaHelper(db database.Connection, configID int64) *dbNodemetaHelper {
	return &dbNodemetaHelper{
		BaseHelper: database.NewBaseHelper(db, "tree_fields", func(t *dbNodemeta) map[string]interface{} {
			return map[string]interface{}{
				"id":        &t.ID,
				"parent_id": &t.ParentID,
				"config_id": &t.ConfigID,

				"name":     &t.Name,
				"type":     &t.Type,
				"required": &t.Required,
			}
		}),
		configID: configID,
	}
}

func (h *dbNodemetaHelper) Create(ctx context.Context, t *dbNodemeta) (*dbNodemeta, error) {
	t.ConfigID = h.configID
	return h.BaseHelper.Create(ctx, t)
}

func (h *dbNodemetaHelper) Get(ctx context.Context, fields []string, condition database.Condition[database.MysqlCondition]) ([]dbNodemeta, error) {
	return h.BaseHelper.Get(ctx, fields,
		condition.And(database.NewMysqlConditionHelper().Set("config_id", database.ConditionOperationEqual, h.configID)))
}

func (h *dbNodemetaHelper) Update(ctx context.Context, t *dbNodemeta, fields []string, condition database.Condition[database.MysqlCondition]) error {
	return h.BaseHelper.Update(ctx, t, fields,
		condition.And(database.NewMysqlConditionHelper().Set("config_id", database.ConditionOperationEqual, h.configID)))
}

func (h *dbNodemetaHelper) Delete(ctx context.Context, condition database.Condition[database.MysqlCondition]) error {
	return h.BaseHelper.Delete(ctx,
		condition.And(database.NewMysqlConditionHelper().Set("config_id", database.ConditionOperationEqual, h.configID)))
}

type dbNodeVal struct {
	database.TableID[int64]
	ParentID *int64 `json:"parent_id,omitempty"`
	GroupID  *int64 `json:"group_id,omitempty"`

	NodemetaID int64  `json:"node_meta_id"`
	Value      string `json:"value,omitempty"`
}

func NewDbNodeValHelper(db database.Connection, valueTable string) *database.BaseHelper[dbNodeVal] {
	return database.NewBaseHelper(db, valueTable, func(t *dbNodeVal) map[string]interface{} {
		return map[string]interface{}{
			"id":           &t.ID,
			"parent_id":    &t.ParentID,
			"group_id":     &t.GroupID,
			"node_type_id": &t.NodemetaID,
			"value":        &t.Value,
		}
	})
}
