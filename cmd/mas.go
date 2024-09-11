/*
Copyright © 2024 Charles <AsciiFaceman> Corbett
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// masCmd represents the mas command
var masCmd = &cobra.Command{
	Use:   "mas",
	Short: "Simple ad-hoc scrapes against mastodon",
	Long:  `Simple ad-hoc scrapes against mastodon.`,
}

func init() {
	rootCmd.AddCommand(masCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// masCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// masCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
