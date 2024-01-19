package queries

import (
	"context"

	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/model"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

type UserRoleMapping struct {
	User  model.User            `json:"user"`
	Roles []RoleWithReferenceID `json:"roles"`
}

type RoleWithReferenceID struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ReferenceID *int64 `json:"reference_id,omitempty"`
}

func GetUserRoleMapping(ctx context.Context, db database.Connection) (map[int64]*UserRoleMapping, error) {
	query := `SELECT role.id, role.name, user_role_mapping.reference_id, user_role_mapping.user_id, user.username
		FROM role
			INNER JOIN user_role_mapping ON role.id = user_role_mapping.role_id
			INNER JOIN user ON user_role_mapping.user_id = user.id;`

	type roleWithUserIDandReferenceID struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		ReferenceID *int64 `json:"reference_id"`
		UserID      int64  `json:"user_id"`
		UserName    string `json:"user_name"`
	}

	resp, err := database.QueryScanner(ctx, db, func(r *roleWithUserIDandReferenceID) []interface{} {
		return []interface{}{
			&r.ID,
			&r.Name,
			&r.ReferenceID,
			&r.UserID,
			&r.UserName,
		}
	}, query)
	if err != nil {
		return nil, err
	}

	userRoleMapping := make(map[int64]*UserRoleMapping)
	for _, r := range resp {
		if _, ok := userRoleMapping[r.UserID]; !ok {
			userRoleMapping[r.UserID] = &UserRoleMapping{
				User: model.User{
					TableID:  database.TableID[int64]{ID: r.UserID},
					Username: r.UserName,
				},
			}
		}

		userRoleMapping[r.UserID].Roles = append(userRoleMapping[r.UserID].Roles, RoleWithReferenceID{
			ID:          r.ID,
			Name:        r.Name,
			ReferenceID: r.ReferenceID,
		})
	}

	return userRoleMapping, nil
}

func GetRoleFromUserID(ctx context.Context, db database.Connection, userID int64) ([]RoleWithReferenceID, error) {
	query := `SELECT role.id, role.name, user_role_mapping.reference_id
		FROM role
			INNER JOIN user_role_mapping ON role.id = user_role_mapping.role_id
		WHERE user_role_mapping.user_id = ?;`

	return database.QueryScanner(ctx, db, func(r *RoleWithReferenceID) []interface{} {
		return []interface{}{
			&r.ID,
			&r.Name,
			&r.ReferenceID,
		}
	}, query, userID)
}
