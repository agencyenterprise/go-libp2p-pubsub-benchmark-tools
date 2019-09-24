package main

import (
	"os"

	"github.com/agencyenterprise/gossip-host/internal/config"
	"github.com/agencyenterprise/gossip-host/internal/host"
	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func setup() {
	var (
		confLoc, listens, rpcListen, peers string
	)

	rootCmd = &cobra.Command{
		Use:   "start",
		Short: "Start node",
		Long:  `Starts the gossip pub/sub node`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Infof("Loading config: %s", confLoc)
			conf, err := config.Load(confLoc, listens, rpcListen, peers)
			if err != nil {
				logger.Errorf("error loading config\n%v", err)
				return err
			}
			logger.Info("Loaded configuration. Starting host...")

			if err = host.Start(conf); err != nil {
				logger.Errorf("err starting host\n%v", err)
				return err
			}

			// note: we are capturing the ctrl+c signal and need to exit, here
			os.Exit(0)

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&confLoc, "config", "c", "config.json", "The configuration file.")
	rootCmd.PersistentFlags().StringVarP(&listens, "listens", "l", "", "Addresses on which to listen. Comma separated. Overides config.json.")
	rootCmd.PersistentFlags().StringVarP(&peers, "peers", "p", "", "Peers to connect. Comma separated. Overides config.json.")
	rootCmd.PersistentFlags().StringVarP(&rpcListen, "rpc-listen", "r", "", "RPC listen address. Overides config.json.")
}

func main() {
	logger.Set(logger.ContextHook{}, false)

	setup()

	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("err executing command\n%v", err)
		os.Exit(1)
	}
}
