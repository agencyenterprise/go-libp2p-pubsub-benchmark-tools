package main

import (
	"context"
	"os"

	"github.com/agencyenterprise/gossip-host/pkg/host"
	"github.com/agencyenterprise/gossip-host/pkg/subnet/config"
	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func setup() *cobra.Command {
	var (
		confLoc string
	)

	rootCmd := &cobra.Command{
		Use:   "start",
		Short: "Start subnet",
		Long:  `Start a subnet of interconnected gossipsub hosts`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := logger.Set(logger.ContextHook{}, "", false); err != nil {
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

			// create a context
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// create the host
			// note: performed conf nil check, above
			h, err := host.New(ctx, *conf)
			if err != nil {
				logger.Errorf("err creating new host:\n%v", err)
				return err
			}

			// start the server
			// note: this is blocking
			if err = h.Start(); err != nil {
				logger.Errorf("err starting host\n%v", err)
				return err
			}

			// note: we are capturing the ctrl+c signal and need to exit, here
			os.Exit(0)

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&confLoc, "config", "c", "host.config.json", "The configuration file.")

	return rootCmd
}

func main() {
	rootCmd := setup()

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("err executing command\n%v", err)
	}
}
