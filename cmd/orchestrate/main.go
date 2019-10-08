package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agencyenterprise/gossip-host/pkg/client"
	"github.com/agencyenterprise/gossip-host/pkg/logger"
	"github.com/agencyenterprise/gossip-host/pkg/orchestrate/config"
	"github.com/agencyenterprise/gossip-host/pkg/subnet"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// note: min/max are inclusive
func randBetween(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func setup() *cobra.Command {
	var (
		confLoc, msgLoc, loggerLoc string
	)

	rootCmd := &cobra.Command{
		Use:   "start",
		Short: "Orchestrate clients and optionally a subnet",
		Long:  `Spins up clients and optionally a subnet and sends the hosts messages at the specified interval.`,
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

			var hostAddrs []string
			if !conf.Orchestra.OmitSubnet {
				sConf := config.BuildSubnetConfig(conf)

				// create the subnet
				snet, err := subnet.New(&subnet.Props{
					CTX:  ctx,
					Conf: sConf,
				})
				if err != nil {
					logger.Errorf("err creating subnet:\n%v", err)
					return err
				}

				// start the subnet
				if !conf.Orchestra.OmitSubnet {
					if err = snet.Start(); err != nil {
						logger.Errorf("err starting subnet\n%v", err)
						return err
					}
				}

				hostAddrs = snet.Addresses()
			}

			if conf.Orchestra.OmitSubnet {
				hostAddrs = conf.Orchestra.HostsIfOmitSubnet
			}

			if len(hostAddrs) == 0 {
				logger.Errorf("no host addresses:\n%v", err)
				return err
			}

			ticker := time.NewTicker(time.Duration(conf.Orchestra.MessageNanoSecondInterval) * time.Nanosecond)
			defer ticker.Stop()

			eChan := make(chan error)
			select {
			case <-stop:
				// note: I don't like '^C' showing up on the same line as the next logged line...
				fmt.Println("")
				logger.Info("Received stop signal from os. Shutting down...")

			case <-ticker.C:
				go func() {
					id, err := uuid.NewRandom()
					if err != nil {
						logger.Errorf("err generating uuid:\n%v", err)
						eChan <- err
					}

					peerIdx := randBetween(0, len(hostAddrs)-1)
					peer := hostAddrs[peerIdx]

					logger.Infof("sending message to %s for gossip", peer)
					if err := client.Gossip(id[:], conf.Orchestra.MessageLocation, peer, conf.Orchestra.MessageByteSize, conf.Orchestra.ClientTimeoutSeconds); err != nil {
						logger.Fatalf("err sending messages\n%v", err)
						eChan <- err
					}
				}()

			case err := <-eChan:
				logger.Errorf("received err on channel:\n%v", err)
				return err
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&confLoc, "config", "c", "configs/subnet/config.json", "The configuration file.")
	rootCmd.PersistentFlags().StringVarP(&loggerLoc, "log", "", "", "Log file location. Defaults to standard out.")
	rootCmd.Flags().StringVarP(&msgLoc, "message", "m", "client.message.json", "The message file to send to peers.")

	return rootCmd
}

func main() {
	rootCmd := setup()

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("err executing command\n%v", err)
	}

	logger.Info("done")
}
