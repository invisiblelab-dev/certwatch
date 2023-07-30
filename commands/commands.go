package commands

import (
	"errors"

	"github.com/invisiblelab-dev/certwatch/factory"
	"github.com/spf13/cobra"
)

var ErrSilent = errors.New("SilentErr")

func Parse() *cobra.Command {
	f := &factory.Factory{}
	f.NotifierService = factory.NewNotifierService(f)

	rootCmd := &cobra.Command{
		Use:           "certwatch",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())

		return ErrSilent
	})

	rootCmd.AddCommand(newAddDomainCommand())
	rootCmd.AddCommand(newCheckAllCertificatesCommand(f))
	rootCmd.AddCommand(newCheckCertificatesCommand())

	return rootCmd
}
