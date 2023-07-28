package commands

import (
	"errors"
	"fmt"

	"github.com/invisiblelab-dev/certwatch"
	"github.com/invisiblelab-dev/certwatch/config"
	"github.com/invisiblelab-dev/certwatch/factory"
	"github.com/spf13/cobra"
)

var ErrSilent = errors.New("SilentErr")

func Parse() *cobra.Command {
	var cfg certwatch.Config
	f := &factory.Factory{}
	f.NotifierService = factory.NewNotifierService(f)

	rootCmd := &cobra.Command{
		Use:           "certwatch",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			path, err := cmd.Flags().GetString("config")
			if err == nil {
				cfg, err = config.ReadYaml(path)
				f.Config = &cfg
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}
			}

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
	rootCmd.AddCommand(newCheckAllCertificatesCommand(&cfg))
	rootCmd.AddCommand(newCheckCertificatesCommand())

	return rootCmd
}
