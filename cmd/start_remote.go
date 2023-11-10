package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thalesfsp/mole/core"
)

const (
	RemoteForwardDoc = `Remote Forwarding allows anyone to expose a service running locally to a remote machine.

This could be particular useful for giving someone on the outside access to an internal web application, for example.

Source endpoints are addresses on the jump server where clients can connect to access services running on the corresponding destination endpoints.
Destination endpoints are addresses of services running on the same machine where mole is getting executed.`
)

var startRemoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Starts a ssh remote port forwarding tunnel",
	Long:  fmt.Sprintf("Starts a ssh remote port forwarding tunnel.\n%s", RemoteForwardDoc),
	Args: func(cmd *cobra.Command, args []string) error {
		conf.TunnelType = "remote"
		return nil
	},
	Run: func(cmd *cobra.Command, arg []string) {
		client := core.New(conf)

		err := client.Start()
		if err != nil {
			os.Exit(1)
		}

	},
}

func init() {
	err := bindFlags(conf, startRemoteCmd)
	if err != nil {
		log.WithError(err).Error("error parsing command line arguments")
		os.Exit(1)
	}

	startCmd.AddCommand(startRemoteCmd)
}
