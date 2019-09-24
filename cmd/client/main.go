package main

import (
	"os"

	"github.com/agencyenterprise/gossip-host/internal/client"
	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func setup() {
	var (
		msgLoc, peers string
	)

	rootCmd = &cobra.Command{
		Use:   "start",
		Short: "Start client",
		Long:  `Starts the client and sends the message to the peers`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Infof("sending message to peers")
			if err := client.Send(msgLoc, peers); err != nil {
				logger.Errorf("err sending messages\n%v", err)
				return err
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&msgLoc, "message", "m", "message.json", "The message file to send to peers.")
	rootCmd.PersistentFlags().StringVarP(&peers, "peers", "p", "", "Peers to connect. Comma separated.")
}

func main() {
	logger.Set(logger.ContextHook{}, false)

	setup()

	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("err executing command\n%v", err)
		os.Exit(1)
	}
}
