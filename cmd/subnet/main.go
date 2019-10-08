package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/agencyenterprise/gossip-host/pkg/logger"
	"github.com/agencyenterprise/gossip-host/pkg/subnet"
	"github.com/agencyenterprise/gossip-host/pkg/subnet/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func setup() *cobra.Command {
	var (
		confLoc, loggerLoc string
	)

	rootCmd := &cobra.Command{
		Use:   "start",
		Short: "Start subnet",
		Long:  `Start a subnet of interconnected gossipsub hosts`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := logger.Set(logger.ContextHook{}, loggerLoc, false); err != nil {
				logrus.Errorf("err initiating logger:\n%v", err)
				return err
			}

			logger.Infof("Loading config: %s", confLoc)
			conf, err := config.Load(confLoc)
			if err != nil {
				logger.Errorf("error loading config\n%v", err)
				return err
			}
			logger.Infof("Loaded configuration. Starting host.\n%v", conf)

			// check the logger location in the conf file
			if conf.General.LoggerLocation != "" {
				switch loggerLoc {
				case conf.General.LoggerLocation:
					break

				case "":
					logger.Warnf("logs will now be written to %s", conf.General.LoggerLocation)
					if err = logger.SetLoggerLoc(conf.General.LoggerLocation); err != nil {
						logger.Errorf("err setting log location to %s:\n%v", conf.General.LoggerLocation, err)
						return err
					}

					break

				default:
					logger.Warnf("log location confliction between flag (%s) and config file (%s); defering to flag (%s)", loggerLoc, conf.General.LoggerLocation, loggerLoc)
				}
			}

			// capture the ctrl+c signal
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, syscall.SIGINT)

			// create a context
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// start the subnet
			if err = subnet.Start(ctx, conf); err != nil {
				logger.Errorf("err starting subnet\n%v", err)
				return err
			}

			select {
			case <-stop:
				// note: I don't like '^C' showing up on the same line as the next logged line...
				fmt.Println("")
				logger.Info("Received stop signal from os. Shutting down...")
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&confLoc, "config", "c", "configs/subnet/config.json", "The configuration file.")
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
