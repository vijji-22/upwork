package httphelper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/constant"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/logger"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/userhelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
)

type MiddlewareFuncWithNext func(http.ResponseWriter, *http.Request, func())

type MiddlewareFuncWithAbort MiddlewareFuncWithNext

func CORSMiddleware(origin, header string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", origin) // TODO: change this to specific domain
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", header) // TODO: change this to specific domain
		res.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func LoggerMiddleware(log logger.Logger) MiddlewareFuncWithNext {
	return func(res http.ResponseWriter, req *http.Request, next func()) {
		start := time.Now()
		logwithid := log.WithField(constant.CtxKey_RequestID, utils.GenerateUUID())

		logwithid.Infof("Request received from %v %v %v", req.RemoteAddr, req.Method, req.URL)
		utils.AddValueInRequestContext(req, constant.CtxKey_Logger, logwithid)

		next()

		logwithid.Infof("Request completed in %v", time.Since(start))
	}
}

// RecoveryMiddleware recovers from panic and send error response
func RecoveryMiddleware(log logger.Logger) MiddlewareFuncWithNext {
	return func(res http.ResponseWriter, req *http.Request, next func()) {
		defer func() {
			if err := recover(); err != nil {
				SendError("intern server error", fmt.Errorf("%v", err), res)
				log.Error(fmt.Errorf("%v", err), "Panic recovered")
			}
		}()

		next()
	}
}

// func DatabaseConnectionMiddleware(db *sql.DB) http.HandlerFunc {
// 	return func(res http.ResponseWriter, req *http.Request) {
// 		utils.AddValueInRequestContext(req, constant.CtxKey_DbConnection, db)
// 	}
// }

func JWTAuthMiddleware[USER database.TableWithID[IDTYPE], IDTYPE int64 | string](loginRequired bool, secret string) MiddlewareFuncWithAbort {
	return func(res http.ResponseWriter, req *http.Request, abort func()) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" && !loginRequired {
			return
		}

		user, reason, err := userhelper.JWTAuthValidate[USER](authHeader, secret)
		if err != nil {
			SendError("authentication failed", err, res)
			abort()
			return
		}

		if reason != "login" {
			SendError("authentication failed", fmt.Errorf("some issue with token"), res)
			abort()
			return
		}

		utils.AddValueInRequestContext(req, constant.CtxKey_User, user)
	}
}
