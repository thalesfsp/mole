package cmd

import (
	"errors"
	"os"

	"github.com/thalesfsp/mole/alias"
	"github.com/thalesfsp/mole/core"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var startAliasCmd = &cobra.Command{
	Use:   "alias [name]",
	Short: "Starts a ssh tunnel by alias",
	Long: `Starts a ssh tunnel by alias

The flags provided through this command can be used to override the one with the
same name stored in the alias.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("alias name not provided")
		}

		aliasName = args[0]

		return nil
	},
	Run: func(cmd *cobra.Command, arg []string) {
		var err error

		// This can't be inside init() because of https://github.com/spf13/cobra/issues/1019
		cmd.Flags().VisitAll(func(f *flag.Flag) {
			if f.Changed {
				givenFlags = append(givenFlags, f.Name)
			}
		})

		al, err := alias.Get(aliasName)
		if err != nil {
			log.WithError(err).Errorf("failed to start tunnel from alias %s", aliasName)
			os.Exit(1)
		}

		err = conf.Merge(al, givenFlags)
		if err != nil {
			log.WithError(err).Errorf("failed to start tunnel from alias %s", aliasName)
			os.Exit(1)
		}

		client := core.New(conf)

		err = client.Start()
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"alias": aliasName,
			}).Errorf("failed to start tunnel from alias %s", aliasName)
			os.Exit(1)
		}
	},
}

func init() {
	startAliasCmd.Flags().BoolVarP(&conf.Verbose, "verbose", "v", false, "increase log verbosity")
	startAliasCmd.Flags().BoolVarP(&conf.Insecure, "insecure", "i", false, "skip host key validation when connecting to ssh server")
	startAliasCmd.Flags().BoolVarP(&conf.Detach, "detach", "x", false, "run process in background")

	startCmd.AddCommand(startAliasCmd)
}
