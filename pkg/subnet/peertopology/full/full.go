package full

import (
	"github.com/agencyenterprise/gossip-host/pkg/host"
	"github.com/agencyenterprise/gossip-host/pkg/logger"
)

// Build connects the hosts using the linear topology
func Build(hosts []*host.Host) error {
	var err error

	for i := range hosts {
		for j := range hosts {
			if i == j {
				continue
			}

			if err = hosts[i].Connect(hosts[j].IFPSAddresses()); err != nil {
				logger.Errorf("err connecting %s with %s:\n%v", hosts[i].ID(), hosts[j].ID(), err)
				return err
			}
		}
	}

	return nil
}
