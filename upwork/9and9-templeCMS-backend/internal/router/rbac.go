package router

import (
	"github.com/gin-gonic/gin"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/config"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/handler"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/model"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/ginhelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/httphelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac/mysql"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/userhelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
)

// rbacRoutes will register all routes related to rbac. [access, role, user, role_access_mapping, user_role_mapping, auth/login]
// [POST] 					/api/v1/auth/login
// [GET, POST, PUT, DELETE] /api/v1/access
// [GET, POST, PUT, DELETE] /api/v1/role
// [GET, POST, PUT, DELETE] /api/v1/user
// [GET, POST, PUT, DELETE] /api/v1/role_access_mapping
// [GET, POST, PUT, DELETE] /api/v1/user_role_mapping
func rbacRoutes(v1 *gin.RouterGroup, db database.Connection, conf *config.Config) {

	// [GET, POST, PUT, DELETE] /api/v1/access
	ginhelper.Register(v1, httphelper.NewCrudHelper(mysql.NewAccessHelper(db), utils.ParseInt, EmptyCondition))

	// [GET, POST, PUT, DELETE] /api/v1/role
	ginhelper.Register(v1, httphelper.NewCrudHelper(mysql.NewRoleHelper(db), utils.ParseInt, EmptyCondition))

	// [GET, POST, PUT, DELETE] /api/v1/role_access_mapping
	ginhelper.Register(v1, httphelper.NewCrudHelper(mysql.NewRoleAccessMappingHelper(db), utils.ParseInt, EmptyCondition))

	userHelper := userhelper.NewUserHelper(model.NewUserHelper(db), conf.App.JWTSecret, model.NewUserValidator())
	// [GET, POST, PUT, DELETE] /api/v1/user
	v1.Group("auth").POST("/login", ginhelper.HttpHandlerToGinHandler(httphelper.LoginHandler(userHelper, EmptyCondition)))

	// [GET, POST, PUT, DELETE] /api/v1/user_role_mapping
	ginhelper.Register(v1, httphelper.NewCrudHelper(userHelper, utils.ParseInt, EmptyCondition))

	// [GET, POST, PUT, DELETE] /api/v1/user_role_mapping
	ginhelper.Register(v1, httphelper.NewCrudHelper(mysql.NewUserRoleMappingHelper(db), utils.ParseInt, EmptyCondition))

	v1.GET("/user_role_mapping/detail", handler.GetUserRoleMapping(db))
	v1.GET("/user_role_mapping/detail/:user_id", handler.GetRoleFromUserID(db))
}
