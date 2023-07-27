package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

var ErrSilent = errors.New("SilentErr")

func Parse() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "certwatch",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())

		return ErrSilent
	})

	rootCmd.AddCommand(newAddDomainCommand())
	rootCmd.AddCommand(newCheckAllCertificatesCommand())
	rootCmd.AddCommand(newCheckCertificatesCommand())

	return rootCmd
}
