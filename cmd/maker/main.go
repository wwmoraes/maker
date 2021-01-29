package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wwmoraes/maker"
)

var rootCmd = &cobra.Command{
	Use:               "maker",
	Short:             "make-all-the-things Makefile includes manager",
	Long:              "manages composable make snippets that can be included on a main Makefile, and are configurable through variables and chain rules",
	PersistentPreRunE: preRun,
	Args:              cobra.NoArgs,
	SilenceUsage:      true,
}

var (
	mk *maker.Maker
)

func preRun(cmd *cobra.Command, args []string) (err error) {
	mk, err = maker.NewDefault()

	return err
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
