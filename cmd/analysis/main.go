package main

import (
	"encoding/json"

	"github.com/agencyenterprise/gossip-host/pkg/analysis"
	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func setup() *cobra.Command {
	var (
		loggerLoc string
	)

	rootCmd := &cobra.Command{
		Use:   "analyze [log file]",
		Short: "Analyze a log file",
		Long:  `Analyzes a log file and outputs the metrics to standard out or the specified log file`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := logger.Set(logger.ContextHook{}, loggerLoc, false); err != nil {
				logrus.Errorf("err initiating logger:\n%v", err)
				return err
			}

			logger.Infof("analyzing log file at %s", args[0])
			metrics, err := analysis.Analyze(args[0])
			if err != nil {
				logger.Errorf("err analyzing log file %s:\n%v", args[1], err)
				return err
			}

			for _, metric := range metrics {
				js, err := json.Marshal(metric)
				if err != nil {
					logger.Errorf("err marshalling metric json:\n%v", err)
					return err
				}
				logger.Info(string(js))
			}

			return nil
		},
	}

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
