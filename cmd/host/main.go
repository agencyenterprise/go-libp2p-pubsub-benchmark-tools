package main

import (
	"os"

	"github.com/agencyenterprise/gossip-host/internal/config"
	"github.com/agencyenterprise/gossip-host/internal/host"
	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func setup() *cobra.Command {
	var (
		confLoc, listens, rpcListen, peers, loggerLoc string
	)

	rootCmd := &cobra.Command{
		Use:   "start",
		Short: "Start node",
		Long:  `Starts the gossip pub/sub node`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := logger.Set(logger.ContextHook{}, loggerLoc, false); err != nil {
				logrus.Errorf("err initiating logger:\n%v", err)
				return err
			}

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
	rootCmd.PersistentFlags().StringVarP(&loggerLoc, "log", "", "", "Log file location. Defaults to standard out.")

	return rootCmd
}

func main() {
	rootCmd := setup()

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("err executing command\n%v", err)
	}
}
