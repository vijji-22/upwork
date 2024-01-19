package router

import (
	"github.com/gin-gonic/gin"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/config"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/handler"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/model"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/ginhelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/httphelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/logger"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/tree"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
)

type gitRouterHandler struct {
	log  logger.Logger
	conf *config.Config
	db   database.Connection
}

func NewRouterHandler(log logger.Logger, conf *config.Config, db database.Connection) *gitRouterHandler {
	return &gitRouterHandler{
		log:  log,
		conf: conf,
		db:   db,
	}
}

func (h *gitRouterHandler) RegisterRoute(r *gin.Engine) {
	v1 := r.Group("api/v1")

	authMiddleware := ginhelper.HttpHandlerToGinHandlerWithAbort(httphelper.JWTAuthMiddleware[model.User](true, h.conf.App.JWTSecret))

	ginhelper.Register(v1, httphelper.NewCrudHelper(model.GetTemplateHelper(h.db), utils.ParseInt, templateConditionFactory), authMiddleware)
	ginhelper.Register(v1, httphelper.NewCrudHelper(model.GetTempleHelper(h.db), utils.ParseInt, templeConditionFactory), authMiddleware)

	ginhelper.Register(v1, httphelper.NewCrudHelper(model.GetServiceHelper(h.db), utils.ParseInt, serviceConditionFactory), authMiddleware)
	ginhelper.Register(v1, httphelper.NewCrudHelper(model.GetServiceTypeHelper(h.db), utils.ParseInt, serviceTypeConditionFactory), authMiddleware)

	// temple config meta and value APIs
	ginhelper.Register(v1, httphelper.NewCrudHelper(tree.NewDbNodemetaHelper(h.db, 1),
		utils.ParseInt, EmptyCondition).WithModelName("temple_config_meta"))
	v1.GET("temple_config_meta/tree", handler.GetMetaTreeHandler(h.db, 1))

	ginhelper.Register(v1, httphelper.NewCrudHelper(tree.NewDbNodeValHelper(h.db, "temple_config"),
		utils.ParseInt, EmptyCondition).WithModelName("temple_config_value"))

	v1.POST("temple_config_value/tree/init/:temple_id", handler.InitNewTempleConfigFromMetaHandler(h.db, 1))
	v1.GET("temple_config_value/tree/:temple_id", handler.GetTempleConfigHandler(h.db, 1))

	// register all rbac routes
	rbacRoutes(v1, h.db, h.conf)
}
