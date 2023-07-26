package commands

import (
	"github.com/invisiblelab-dev/certwatch"
	"github.com/invisiblelab-dev/certwatch/runners"
	"github.com/spf13/cobra"
)

func newCheckCertificatesCommand() *cobra.Command {
	opts := certwatch.CheckCertificatesOptions{}
	checkCertificatesCommand := &cobra.Command{
		Use:   "check",
		Short: "Check if given domains are close to end",
		Run: func(cmd *cobra.Command, args []string) {
			runners.RunCheckCertificatesCommand(opts)
		},
	}
	checkCertificatesCommand.Flags().StringSliceVar(&opts.Domains, "domain", []string{}, "domains to check, separated by comma")
	return checkCertificatesCommand
}

func newCheckAllCertificatesCommand() *cobra.Command {
	opts := certwatch.CheckAllCertificatesOptions{}
	checkCertificatesCommand := &cobra.Command{
		Use:   "check-all",
		Short: "Check if your added domains are close to end",
		Run: func(cmd *cobra.Command, args []string) {
			runners.RunCheckAllCertificatesCommand(opts)
		},
	}

	checkCertificatesCommand.Flags().BoolVar(&opts.Force, "force", false, "force check every domain")

	return checkCertificatesCommand
}
