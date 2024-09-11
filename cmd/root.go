/*
Copyright © 2024 Charles <AsciiFaceman> Corbett
*/
package cmd

import (
	"os"

	"github.com/asciifaceman/emri/pkg/global"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "emri",
	Short: "Emri is a tool for documenting the fediverse.",
	Long:  `Emri is a tool for documenting the fediverse..`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.emri.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initialize global config and logging
func init() {
	if err := global.InitLogging(zap.DebugLevel); err != nil {
		zap.S().Fatalw("failed to init root logger", "err", err)
	}

	if err := global.InitConfig(cfgFile); err != nil {
		zap.S().Fatalw("failed to load config file", "err", err)
	}

	var initLevel zapcore.Level
	if global.C().Verbose {
		initLevel = zap.DebugLevel
	} else {
		initLevel = zap.InfoLevel
	}
	if err := global.InitLogging(initLevel); err != nil {
		zap.S().Fatalw("failed to init root logger", "err", err)
	}
}
