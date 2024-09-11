/*
Copyright © 2024 Charles <AsciiFaceman> Corbett
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "emri API",
	Long:  `emri API.`,
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
