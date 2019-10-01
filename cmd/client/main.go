package main

import (
	"os"

	"github.com/agencyenterprise/gossip-host/pkg/client"
	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func setup() *cobra.Command {
	var (
		msgLoc, peers, loggerLoc string
	)

	rootCmd := &cobra.Command{
		Use:   "start",
		Short: "Start client",
		Long:  `Starts the client and sends the message to the peers`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := logger.Set(logger.ContextHook{}, loggerLoc, false); err != nil {
				logrus.Errorf("err initiating logger:\n%v", err)
				return err
			}

			logger.Infof("sending message to peers")
			if err := client.Send(msgLoc, peers); err != nil {
				logger.Errorf("err sending messages\n%v", err)
				return err
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&msgLoc, "message", "m", "client.message.json", "The message file to send to peers.")
	rootCmd.PersistentFlags().StringVarP(&peers, "peers", "p", ":8080", "Peers to connect. Comma separated.")
	rootCmd.PersistentFlags().StringVarP(&loggerLoc, "log", "", "", "Log file location. Defaults to standard out.")

	return rootCmd
}

func main() {
	rootCmd := setup()

	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("err executing command\n%v", err)
		os.Exit(1)
	}
}
