package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thalesfsp/mole/core"
)

var (
	follow      bool
	showLogsCmd = &cobra.Command{
		Use:   "logs [name]",
		Short: "Shows log messages of a detached running application instance",
		Long:  "Shows log messages of a detached running application instance",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				id = args[0]
			}

			return nil
		},
		Run: func(cmd *cobra.Command, arg []string) {
			err := core.ShowLogs(id, follow)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"id": id,
				}).Error("error opening log file")
				os.Exit(1)
			}
		},
	}
)

func init() {
	showLogsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "output appended data as the file grows")
	showCmd.AddCommand(showLogsCmd)
}
