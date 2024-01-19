package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/queries"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/httphelper"
)

func GetUserRoleMapping(db database.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := queries.GetUserRoleMapping(c, db)
		if err != nil {
			httphelper.SendError("Error while fetching user role mapping", err, c.Writer)
			return
		}

		httphelper.SendData(resp, c.Writer)
	}
}

func GetRoleFromUserID(db database.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, found := c.Params.Get("user_id")
		if !found {
			httphelper.SendError("User id not found", nil, c.Writer)
			return
		}

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			httphelper.SendError("Invalid user id", err, c.Writer)
			return
		}

		resp, err := queries.GetRoleFromUserID(c, db, int64(userIDInt))
		if err != nil {
			httphelper.SendError("Error while fetching role for user", err, c.Writer)
			return
		}

		httphelper.SendData(resp, c.Writer)
	}
}
