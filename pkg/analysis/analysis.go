package analysis

import (
	"github.com/agencyenterprise/gossip-host/pkg/analysis/types"
	"github.com/agencyenterprise/gossip-host/pkg/logger"
)

// Analyze parses a log file and returns performance metrics
func Analyze(logLoc string) ([]*types.Metric, error) {
	scanner, logFileCloser, err := loadLogFile(logLoc)
	if err != nil {
		logger.Errorf("err opening log file at %s:\n%v", logLoc, err)
		return nil, err
	}
	defer logFileCloser()

	return buildMetricsFromScanner(scanner)
}
