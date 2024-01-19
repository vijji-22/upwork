package model

import (
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac/mysql"
)

type Temple struct {
	database.TableID[int64]
	Name      string `json:"name,omitempty"`
	Logo      string `json:"logo,omitempty"`
	HeroImage string `json:"heroimage,omitempty"`

	Timings Timings `json:"timings,omitempty"`

	Address Address `json:"address,omitempty"`

	// TODO : adopted temple and AWS config
}

type Timings struct {
	OpenTime  string `json:"opentime,omitempty"`
	CloseTime string `json:"closetime,omitempty"`
}

type Address struct {
	Line1   string `json:"line1,omitempty"`
	Line2   string `json:"line2,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zip     string `json:"zip,omitempty"`
	Country string `json:"country,omitempty"`
}

func GetTempleHelper(db database.Connection) database.CrudHelper[database.MysqlCondition, Temple, int64] {
	h := database.NewBaseHelper(db, "temple", func(t *Temple) map[string]interface{} {
		return map[string]interface{}{
			"id":   &t.ID,
			"name": &t.Name,

			"logo":      &t.Logo,
			"heroimage": &t.HeroImage,

			"opentime":  &t.Timings.OpenTime,
			"closetime": &t.Timings.CloseTime,

			"address_line1": &t.Address.Line1,
			"address_line2": &t.Address.Line2,
			"city":          &t.Address.City,
			"state":         &t.Address.State,
			"zip":           &t.Address.Zip,
			"country":       &t.Address.Country,
		}
	})

	rbacHelper := rbac.NewRbacHelper(mysql.MysqlRbacHelper(db))
	rbacCrudHelper := rbac.NewCrudHelper(rbacHelper, h, UserIDFromCTX)

	return rbacCrudHelper
}
