package commands

import (
	"fmt"

	"github.com/invisiblelab-dev/certwatch"
	"github.com/invisiblelab-dev/certwatch/config"
	"github.com/invisiblelab-dev/certwatch/factory"
	"github.com/invisiblelab-dev/certwatch/runners"
	"github.com/spf13/cobra"
)

func newCheckCertificatesCommand() *cobra.Command {
	opts := certwatch.CheckCertificatesOptions{}
	checkCertificatesCommand := &cobra.Command{
		Use:   "check",
		Short: "Check if given domains are close to end",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runners.RunCheckCertificatesCommand(opts); err != nil {
				return fmt.Errorf("check: %w", err)
			}

			return nil
		},
	}
	checkCertificatesCommand.Flags().StringSliceVar(&opts.Domains, "domain", []string{}, "domains to check, separated by comma")

	return checkCertificatesCommand
}

func newCheckAllCertificatesCommand(f *factory.Factory) *cobra.Command {
	opts := certwatch.CheckAllCertificatesOptions{}
	checkCertificatesCommand := &cobra.Command{
		Use:   "check-all",
		Short: "Check if your added domains are close to end",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.ReadYaml(opts.Path)
			f.Config = &cfg
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if err := runners.RunCheckAllCertificatesCommand(f, &cfg); err != nil {
				return fmt.Errorf("check-all: %w", err)
			}

			return nil
		},
	}

	checkCertificatesCommand.Flags().BoolVar(&opts.Force, "force", false, "force check every domain")
	checkCertificatesCommand.Flags().StringVar(&opts.Path, "config", "certwatch.yaml", "define path to config file")

	return checkCertificatesCommand
}
