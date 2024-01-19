package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/httphelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/tree"
)

func GetMetaTreeHandler(db database.Connection, configID int64) gin.HandlerFunc {
	templeHelper := tree.NewTreeHelper(tree.NewTreeDbHelper(db, "temple_config", configID))

	return func(c *gin.Context) {
		meta, err := templeHelper.GetMeta(c)
		if err != nil {
			httphelper.SendError("Error while fetching meta", err, c.Writer)
			return
		}

		httphelper.SendData(meta, c.Writer)
	}
}

func GetTempleConfigHandler(db database.Connection, configID int64) gin.HandlerFunc {
	templeHelper := tree.NewTreeHelper(tree.NewTreeDbHelper(db, "temple_config", configID))

	return func(c *gin.Context) {

		templeID, found := c.Params.Get("temple_id")
		if !found {
			httphelper.SendError("Temple id not found", nil, c.Writer)
			return
		}

		templeIDInt, err := strconv.Atoi(templeID)
		if err != nil {
			httphelper.SendError("Invalid temple id", err, c.Writer)
			return
		}

		templeTree, err := templeHelper.GetTree(c, int64(templeIDInt))
		if err != nil {
			httphelper.SendError("Error while fetching temple tree", err, c.Writer)
			return
		}

		httphelper.SendData(templeTree, c.Writer)
	}
}

func InitNewTempleConfigFromMetaHandler(db database.Connection, configID int64) gin.HandlerFunc {
	templeHelper := tree.NewTreeHelper(tree.NewTreeDbHelper(db, "temple_config", configID))

	return func(c *gin.Context) {

		templeID, found := c.Params.Get("temple_id")
		if !found {
			httphelper.SendError("Temple id not found", nil, c.Writer)
			return
		}

		templeIDInt, err := strconv.Atoi(templeID)
		if err != nil {
			httphelper.SendError("Invalid temple id", err, c.Writer)
			return
		}

		err = templeHelper.SetupValueFromMeta(c, int64(templeIDInt))
		if err != nil {
			httphelper.SendError("Error while setting up value from meta", err, c.Writer)
			return
		}

		templeTree, err := templeHelper.GetTree(c, int64(templeIDInt))
		if err != nil {
			httphelper.SendError("Error while fetching temple tree", err, c.Writer)
			return
		}

		httphelper.SendData(templeTree, c.Writer)
	}
}
