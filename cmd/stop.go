package cmd

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thalesfsp/mole/core"
)

var (
	stopCmd = &cobra.Command{
		Use:   "stop [alias name or id]",
		Short: "Stops an instance of mole ",
		Long:  "Stops an instance of mole by either a given auto generated id or alias",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("alias name or id not provided")
			}

			conf.Id = args[0]

			return nil
		},
		Run: func(cmd *cobra.Command, arg []string) {
			c := core.New(conf)

			err := c.Stop()
			if err != nil {
				log.WithError(err).Error("error stopping detached mole instance")
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(stopCmd)
}
