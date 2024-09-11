/*
Copyright © 2024 Charles <AsciiFaceman> Corbett
*/
package cmd

import (
	"fmt"

	"github.com/asciifaceman/emri/pkg/dal"
	"github.com/asciifaceman/emri/pkg/dal/models"
	"github.com/asciifaceman/emri/pkg/global"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize emri db (destructive)",
	Long:  `Initialize emri db (destructive).`,
	Run: func(cmd *cobra.Command, args []string) {
		dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=%s",
			global.C().PostgresConfig.Hostname,
			global.C().PostgresConfig.Manage.Username,
			global.C().PostgresConfig.Manage.Password,
			global.C().PostgresConfig.Port,
			global.C().PostgresConfig.SSL,
		)
		dbn := global.C().PostgresConfig.Database
		owner := global.C().PostgresConfig.Runtime.Username
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: dal.NewLogger("gorm", dbn),
		})
		if err != nil {
			zap.S().Fatalw("error connecting to database", "error", err)
		}

		// Check for DB

		stmt := "SELECT * FROM pg_database WHERE datname = '%s';"
		rs := db.Raw(fmt.Sprintf(stmt, dbn))
		if rs.Error != nil {
			zap.S().Fatalw("error checking if db exists", "error", rs.Error)
		}

		// DROP DB

		stmt = "DROP DATABASE %s;"
		var recv = make(map[string]interface{})
		if rs.Find(recv); len(recv) > 0 {
			if cs := db.Exec(fmt.Sprintf(stmt, dbn)); cs.Error != nil {
				zap.S().Fatalw("error dropping database", "database", dbn, "error", cs.Error)
			}
		} else {
			zap.S().Errorw("database doesn't exist, continuing...", "database", dbn)
		}

		// CREATE DB

		stmt = "CREATE DATABASE %s WITH OWNER = %s ENCODING = 'UTF8' IS_TEMPLATE = False;"

		if cs := db.Exec(fmt.Sprintf(stmt, dbn, owner)); cs.Error != nil {
			zap.S().Fatalw("Failed to create database", "error", cs.Error, "database", dbn)
		}

		stmt = "GRANT ALL PRIVILEGES ON DATABASE %s TO %s"

		if gs := db.Exec(fmt.Sprintf(stmt, dbn, owner)); gs.Error != nil {
			zap.S().Fatalw("failed to grant privileges", "error", gs.Error)
		}

		d, err := db.DB()
		if err != nil {
			zap.S().Fatalw("failed to get DB() instance to close", "error", err)
		}

		if err := d.Close(); err != nil {
			zap.S().Fatalw("failed to close connection to reconnect under database", "error", err)
		}

		zap.S().Info("Connecting to db under management creds to grant")

		dsn = fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=%s dbname=%s",
			global.C().PostgresConfig.Hostname,
			global.C().PostgresConfig.Manage.Username,
			global.C().PostgresConfig.Manage.Password,
			global.C().PostgresConfig.Port,
			global.C().PostgresConfig.SSL,
			dbn,
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: dal.NewLogger("gorm", dbn),
		})
		if err != nil {
			zap.S().Fatalw("error connecting to database", "error", err)
		}

		stmt = "GRANT ALL PRIVILEGES ON SCHEMA public TO %s;"
		if gs := db.Exec(fmt.Sprintf(stmt, owner)); gs.Error != nil {
			zap.S().Fatalw("failed to grant privileges on schema", "error", gs.Error)
		}

		zap.S().Info("Connecting to db under app creds to migrate")

		d, err = db.DB()
		if err != nil {
			zap.S().Fatalw("failed to get DB() instance to close", "error", err)
		}

		if err := d.Close(); err != nil {
			zap.S().Fatalw("failed to close connection to reconnect under database", "error", err)
		}

		dl, err := dal.New()
		if err != nil {
			zap.S().Fatalw("failed to get dal", "error", err)
		}

		if err := dl.Migrate(models.Migrate...); err != nil {
			zap.S().Fatalw("failed to migrate", "error", err)
		}

		zap.S().Info("complete!")

	},
}

func init() {
	dbCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
