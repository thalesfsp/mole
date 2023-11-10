package cmd

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/thalesfsp/mole/core"
)

var (
	aliasName  string
	id         string
	conf       = &core.Configuration{}
	givenFlags []string

	rootCmd = &cobra.Command{
		Use:  "mole",
		Long: "Tool to create ssh tunnels focused on resiliency and user experience.",
	}
)

// Execute executes the root command
func Execute() error {
	log.SetOutput(os.Stdout)

	return rootCmd.Execute()
}

func bindFlags(conf *core.Configuration, cmd *cobra.Command) error {
	cmd.Flags().BoolVarP(&conf.Verbose, "verbose", "v", false, "increase log verbosity")
	cmd.Flags().BoolVarP(&conf.Insecure, "insecure", "i", false, "skip host key validation when connecting to ssh server")
	cmd.Flags().BoolVarP(&conf.Detach, "detach", "x", false, "run process in background")
	cmd.Flags().VarP(&conf.Source, "source", "S", `set source endpoint address: [<host>]:<port>
multiple -source conf can be provided`)
	cmd.Flags().VarP(&conf.Destination, "destination", "d", `set destination endpoint address: [<host>]:<port>
multiple -destination conf can be provided`)
	cmd.Flags().VarP(&conf.Server, "server", "s", "set server address: [<user>@]<host>[:<port>]")
	cmd.Flags().StringVarP(&conf.Key, "key", "k", "", "set server authentication key file path")
	cmd.Flags().StringVar(&conf.KeyValue, "key-value", "", "set server authentication key")
	cmd.Flags().DurationVarP(&conf.KeepAliveInterval, "keep-alive-interval", "K", 10*time.Second, "time interval for keep alive packets to be sent")
	cmd.Flags().IntVarP(&conf.ConnectionRetries, "connection-retries", "R", 3, `maximum number of connection retries to the ssh server
provide 0 to never give up or a negative number to disable`)
	cmd.Flags().StringVarP(&conf.SshConfig, "config", "c", "$HOME/.ssh/config", "set config file path")
	cmd.Flags().DurationVarP(&conf.WaitAndRetry, "retry-wait", "w", 3*time.Second, "time to wait before trying to reconnect to ssh server")
	cmd.Flags().StringVarP(&conf.SshAgent, "ssh-agent", "A", "", "unix socket to communicate with a ssh agent")
	cmd.Flags().DurationVarP(&conf.Timeout, "timeout", "t", 3*time.Second, "ssh server connection timeout")
	cmd.Flags().BoolVarP(&conf.Rpc, "rpc", "", false, "enable the rpc server")
	cmd.Flags().StringVarP(&conf.RpcAddress, "rpc-address", "", "127.0.0.1:0", `set the network address of the rpc server.
The default value uses a random free port to listen for requests.
The full address is kept on $HOME/.mole/<id>.`)

	// id is a hidden flag used to carry the unique identifier of the instance to
	// the child process when the `--detached` flag is used.
	cmd.Flags().StringVarP(&conf.Id, core.IdFlagName, "", "", "")
	err := cmd.Flags().MarkHidden(core.IdFlagName)
	if err != nil {
		return err
	}

	err = cmd.MarkFlagRequired("server")
	if err != nil {
		return err
	}

	flag.Visit(func(f *flag.Flag) {
		givenFlags = append(givenFlags, f.Name)
	})

	return nil
}
