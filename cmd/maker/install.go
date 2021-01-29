package main

import (
	"github.com/spf13/cobra"
)

var (
	installCmd = &cobra.Command{
		Use:   "install",
		Short: "installs all snippets",
		Long:  "fetches all snippet files listed as dependency",
		RunE:  updateRun,
	}
	installForce bool
)

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&installForce, "force", "f", false, "ignores the lock file, and checks the files directly")
}

func updateRun(cmd *cobra.Command, args []string) (err error) {
	err = mk.Install(installForce)
	if err != nil {
		return err
	}

	return nil
}
