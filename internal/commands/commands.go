package commands

import (
	"github.com/spf13/cobra"
)

func Parse() *cobra.Command {
	rootCmd := &cobra.Command{Use: "certwatch"}
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(newAddDomainCommand())

	return rootCmd
}
