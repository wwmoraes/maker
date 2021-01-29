package main

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a snippet",
	Long:  "fetches a snippet file and adds it as a dependency",
	RunE:  addRun,
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addRun(cmd *cobra.Command, args []string) (err error) {
	err = mk.Add(args[0])
	if err != nil {
		return err
	}

	return nil
}
