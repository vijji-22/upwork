package main

import (
	"errors"
	"fmt"

	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/config"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/router"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/ginhelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/logger"
	"github.com/spf13/cobra"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var rootCMD = &cobra.Command{}
var log = logger.New()

func init() {

	runServerCMD := &cobra.Command{
		Use: "runserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := config.NewConfigFromEnv()

			if conf.App.JWTSecret == "" {
				return errors.New("JWT secrete is required")
			}

			db := database.Connect(log, conf.DatabaseConfig)
			defer db.Close()

			fmt.Println("Running migration...")
			err := database.RunMigration(db, conf.DatabaseConfig)
			if err != nil {
				return err
			}
			fmt.Println("Migration completed successfully")

			ginRouterHelper := router.NewRouterHandler(log, conf, db)
			log.Info("Starting server...")
			ginhelper.StartServer(log, db, conf.App.Port, ginRouterHelper.RegisterRoute)

			return nil
		},
	}

	rootCMD.AddCommand(runServerCMD)
}

func main() {
	if err := rootCMD.Execute(); err != nil {
		log.Fatal(err, "Failed to execute command")
	}
}
