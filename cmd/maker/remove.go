package main

import (
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "removes a snippet",
	Long:  "deletes the snippet file and removes it as a dependency",
	RunE:  removeRun,
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func removeRun(cmd *cobra.Command, args []string) (err error) {
	err = mk.Remove(args[0])
	if err != nil {
		return err
	}

	return nil
}
