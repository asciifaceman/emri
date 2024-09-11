/*
Copyright © 2024 Charles <AsciiFaceman> Corbett
*/
package cmd

import (
	"github.com/asciifaceman/emri/pkg/social"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// scrapeCmd represents the scrape command
var masScrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape about/blocks/peered from the given domain",
	Long:  `Scrape about/blocks/peered from the given domain.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, domain := range args {
			zap.S().Infow("scraping domain", "domain", domain)

			err := social.Validate(domain)
			if err != nil {
				zap.S().Errorw("error validating domain", "error", err, "domain", domain)
				return
			}

			about, err := social.About(domain)
			if err != nil {
				zap.S().Errorw("error scraping about", "error", err, "domain", domain)
			}
			spew.Dump(about)
		}
	},
}

func init() {
	masCmd.AddCommand(masScrapeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scrapeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scrapeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
