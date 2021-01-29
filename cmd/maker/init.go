package main

import (
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize maker on a directory",
	Long:  "creates an initial maker configuration with the default repository",
	RunE:  initRun,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initRun(cmd *cobra.Command, args []string) error {
	return mk.Init()
}
