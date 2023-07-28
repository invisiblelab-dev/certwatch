package commands

import (
	"fmt"

	"github.com/invisiblelab-dev/certwatch"
	"github.com/invisiblelab-dev/certwatch/runners"
	"github.com/spf13/cobra"
)

func newAddDomainCommand() *cobra.Command {
	opts := certwatch.AddDomainOptions{}
	addDomainCommand := &cobra.Command{
		Use:   "add-domain",
		Short: "Add a new domain and the number of days you want to be notified before it expires",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runners.RunAddDomainCommand(opts); err != nil {
				return fmt.Errorf("add-domain: %w", err)
			}

			return nil
		},
	}

	addDomainCommand.Flags().StringVar(&opts.Domain, "domain", "", "domain to be tracked")
	addDomainCommand.Flags().Int32Var(&opts.DaysBefore, "days", 10, "number of days before expire")
	addDomainCommand.Flags().StringVar(&opts.Path, "path", "certwatch.yaml", "define path to config file")

	cobra.CheckErr(addDomainCommand.MarkFlagRequired("domain"))

	return addDomainCommand
}
