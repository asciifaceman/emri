/*
Copyright © 2024 Charles <AsciiFaceman> Corbett
*/
package cmd

import (
	"github.com/asciifaceman/emri/pkg/dal"
	"github.com/asciifaceman/emri/pkg/dal/models"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// initCmd represents the init command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate emri db (destructive)",
	Long:  `migrate emri db (destructive).`,
	Run: func(cmd *cobra.Command, args []string) {
		zap.S().Info("Connecting to db under app creds to migrate")

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
	dbCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
