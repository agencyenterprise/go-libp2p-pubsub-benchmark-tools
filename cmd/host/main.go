package main

import (
	"context"

	"github.com/agencyenterprise/gossip-host/pkg/host"
	"github.com/agencyenterprise/gossip-host/pkg/host/config"
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
			logger.Infof("Loaded configuration. Starting host.\n%v", conf)

			// check the logger location in the conf file
			if conf.Host.LoggerLocation != "" {
				switch loggerLoc {
				case conf.Host.LoggerLocation:
					break

				case "":
					logger.Warnf("logs will now be written to %s", conf.Host.LoggerLocation)
					if err = logger.SetLoggerLoc(conf.Host.LoggerLocation); err != nil {
						logger.Errorf("err setting log location to %s:\n%v", conf.Host.LoggerLocation, err)
						return err
					}

					break

				default:
					logger.Warnf("log location confliction between flag (%s) and config file (%s); defering to flag (%s)", loggerLoc, conf.Host.LoggerLocation, loggerLoc)
				}
			}

			// create a context
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// create the host
			h, err := host.New(ctx, conf)
			if err != nil {
				logger.Errorf("err creating new host:\n%v", err)
				return err
			}

			// build pubsub
			ps, err := h.BuildPubSub()
			if err != nil || ps == nil {
				logger.Errorf("err building pubsub:\n%v", err)
				return err
			}

			// build rpc
			ch, err := h.BuildRPC(ps)
			if err != nil {
				logger.Errorf("err building rpc:\n%v", err)
				return err
			}

			// connect to peers
			if err = h.Connect(conf.Host.Peers); err != nil {
				logger.Errorf("err connecting to peers:\n%v", err)
				return err
			}

			// start the server
			// note: this is blocking
			if err = h.Start(ch); err != nil {
				logger.Errorf("err starting host\n%v", err)
				return err
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&confLoc, "config", "c", "configs/host/config.json", "The configuration file.")
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

	logger.Info("done")
}
