package mysql

import (
	"context"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

type mysqlRbacHelper struct {
	accessHelper            AccessHelper
	roleHelper              RoleHelper
	roleAccessMappingHelper RoleAccessMappingHelper
	userRoleMappingHelper   UserRoleMappingHelper

	db database.Connection
}

func MysqlRbacHelper(db database.Connection) rbac.DatabaseHelper[database.MysqlCondition, int64] {
	return &mysqlRbacHelper{
		accessHelper:            NewAccessHelper(db),
		roleHelper:              NewRoleHelper(db),
		roleAccessMappingHelper: NewRoleAccessMappingHelper(db),
		userRoleMappingHelper:   NewUserRoleMappingHelper(db),

		db: db,
	}
}

func (mrh *mysqlRbacHelper) GetAccessHelper() database.CrudHelper[database.MysqlCondition, rbac.Access[int64], int64] {
	return mrh.accessHelper
}

func (mrh *mysqlRbacHelper) GetRoleHelper() database.CrudHelper[database.MysqlCondition, rbac.Role[int64], int64] {
	return mrh.roleHelper
}

func (mrh *mysqlRbacHelper) GetRoleAccessMappingHelper() database.CrudHelper[database.MysqlCondition, rbac.RoleAccessMapping[int64], int64] {
	return mrh.roleAccessMappingHelper
}

func (mrh *mysqlRbacHelper) GetUserRoleMappingHelper() database.CrudHelper[database.MysqlCondition, rbac.UserRoleMapping[int64], int64] {
	return mrh.userRoleMappingHelper
}

func (mrh *mysqlRbacHelper) GetAccessForUser(ctx context.Context, userID int64, access string) ([]rbac.AccessWithReferenceID[int64], error) {
	query := `SELECT access.name, access.reference_key, role_access_mapping.project, user_role_mapping.reference_id
		FROM access
			INNER JOIN role_access_mapping ON access.id = role_access_mapping.access_id
			INNER JOIN role ON role_access_mapping.role_id = role.id
			INNER JOIN user_role_mapping ON role.id = user_role_mapping.role_id
		WHERE user_role_mapping.user_id = ? AND access.name = ?`

	return database.QueryScanner(ctx, mrh.db, func(m *rbac.AccessWithReferenceID[int64]) []interface{} {
		return []interface{}{
			&m.AccessName,
			&m.Project,
			&m.ReferenceID,
			&m.ReferenceKey,
		}
	}, query, userID, access)
}

// TODO remove duplicate code
func (mrh *mysqlRbacHelper) GetAccessForUserWithReference(ctx context.Context, userID int64, access string, referenceID int64) ([]rbac.AccessWithReferenceID[int64], error) {
	query := `SELECT access.name, access.reference_key, role_access_mapping.project
		FROM access
			INNER JOIN role_access_mapping ON access.id = role_access_mapping.access_id
			INNER JOIN role ON role_access_mapping.role_id = role.id
			INNER JOIN user_role_mapping ON role.id = user_role_mapping.role_id
		WHERE user_role_mapping.user_id = ? AND access.name = ? AND user_role_mapping.reference_id = ?`

	return database.QueryScanner(ctx, mrh.db, func(m *rbac.AccessWithReferenceID[int64]) []interface{} {
		return []interface{}{
			&m.AccessName,
			&m.Project,
			&m.ReferenceID,
			&m.ReferenceKey,
		}
	}, query, userID, access, referenceID)

}
