/*
Copyright © 2024 Charles <AsciiFaceman> Corbett
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var dbAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add a domain to emri",
	Long:  `add a domain to emri using the same mechanisms as typical usage`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	dbCmd.AddCommand(dbAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
