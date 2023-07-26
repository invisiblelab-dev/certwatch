package commands

import (
	certwatch "github.com/invisiblelab-dev/certwatch"
	"github.com/invisiblelab-dev/certwatch/runners"
	"github.com/spf13/cobra"
)

func newAddDomainCommand() *cobra.Command {
	opts := certwatch.AddDomainOptions{}
	addDomainCommand := &cobra.Command{
		Use:   "add-domain",
		Short: "Add a new domain and the number of days you want to be notified before it expires",
		Run: func(cmd *cobra.Command, args []string) {
			runners.RunAddDomainCommand(opts)
		},
	}

	addDomainCommand.Flags().StringVar(&opts.Domain, "domain", "", "domain to be tracked")
	addDomainCommand.Flags().Int32Var(&opts.DaysBefore, "days", 10, "number of days before expire")
	cobra.CheckErr(addDomainCommand.MarkFlagRequired("domain"))

	return addDomainCommand
}
