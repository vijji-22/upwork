package model

import (
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac/mysql"
)

type Service struct {
	database.TableID[int64]
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	LongDescription  string `json:"long_description"`
	Image            string `json:"image"`

	MaxAmount int64 `json:"max_amount"`
	MinAmount int64 `json:"min_amount"`

	TempleID      int64 `json:"temple_id"`
	ServiceTypeID int64 `json:"service_type_id"`
}

func GetServiceHelper(db database.Connection) database.CrudHelper[database.MysqlCondition, Service, int64] {
	h := database.NewBaseHelper(db, "service", func(s *Service) map[string]interface{} {
		return map[string]interface{}{
			"id":                &s.ID,
			"name":              &s.Name,
			"short_description": &s.ShortDescription,
			"long_description":  &s.LongDescription,
			"image":             &s.Image,
			"max_amount":        &s.MaxAmount,
			"min_amount":        &s.MinAmount,
			"temple_id":         &s.TempleID,
			"service_type_id":   &s.ServiceTypeID,
		}
	})

	rbacHelper := rbac.NewRbacHelper(mysql.MysqlRbacHelper(db))
	return rbac.NewCrudHelper(rbacHelper, h, UserIDFromCTX).WithPublicGet()
}

type ServiceType struct {
	database.TableID[int64]
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	LongDescription  string `json:"long_description"`
	Image            string `json:"image"`
}

func GetServiceTypeHelper(db database.Connection) database.CrudHelper[database.MysqlCondition, ServiceType, int64] {
	h := database.NewBaseHelper(db, "service_type", func(s *ServiceType) map[string]interface{} {
		return map[string]interface{}{
			"id":                &s.ID,
			"name":              &s.Name,
			"short_description": &s.ShortDescription,
			"long_description":  &s.LongDescription,
			"image":             &s.Image,
		}
	})

	rbacHelper := rbac.NewRbacHelper(mysql.MysqlRbacHelper(db))
	return rbac.NewCrudHelper(rbacHelper, h, UserIDFromCTX).WithPublicGet()
}
