package host

import (
	"context"
	"strings"

	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	mplex "github.com/libp2p/go-libp2p-mplex"
	quic "github.com/libp2p/go-libp2p-quic-transport"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	lconfig "github.com/libp2p/go-libp2p/config"
	tcp "github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/multiformats/go-multiaddr"
)

func parseTransportOptions(opts []string) (lconfig.Option, error) {
	var lOpts []lconfig.Option

	for _, opt := range opts {
		switch strings.ToLower(opt) {
		case "tcp":
			lOpts = append(lOpts, libp2p.Transport(tcp.NewTCPTransport))

		case "ws":
			lOpts = append(lOpts, libp2p.Transport(ws.New))

		case "quic":
			lOpts = append(lOpts, libp2p.Transport(quic.NewTransport))

		/* note: utp has a broken gx dep
		case "utp":
			lOpts = append(lOpts, libp2p.Transport(utp.NewUtpTransport))
		*/

		/* note: WIP
		case "udp":
			lOpts = append(lOpts, libp2p.Transport(utp.NewUdpTransport))
		*/

		/* note: need to pass private key? But we didn't for quic...
		case "tls":
			lOpts = append(lOpts, libp2p.Transport(tls.New))
		*/

		case "none":
			if len(opts) > 1 {
				logger.Error("when using the 'none' transport option, cannot also specify other transport options")
				return nil, ErrImproperTransportOption
			}

			return libp2p.NoTransports, nil

		case "default":
			lOpts = append(lOpts, libp2p.DefaultTransports)

		default:
			logger.Errorf("unknown transport option: %s", opt)
			return nil, ErrUnknownTransportOption
		}
	}

	return libp2p.ChainOptions(lOpts...), nil
}

func parseMuxerOptions(opts [][]string) (lconfig.Option, error) {
	var lOpts []lconfig.Option

	for _, opt := range opts {
		if len(opt) != 2 {
			logger.Errorf("improper muxer format, expected ['name', 'type'], received %v", opt)
			return nil, ErrImproperMuxerOption
		}

		switch strings.ToLower(opt[0]) {
		case "yamux":
			lOpts = append(lOpts, libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport))

		case "mplex":
			lOpts = append(lOpts, libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport))

		// TODO: add others?
		default:
			logger.Errorf("unknown muxer option: %s", opt)
			return nil, ErrUnknownMuxerOption
		}
	}

	return libp2p.ChainOptions(lOpts...), nil
}

func parseSecurityOptions(opt string) (lconfig.Option, error) {
	switch strings.ToLower(opt) {
	case "secio":
		return libp2p.Security(secio.ID, secio.New), nil

	case "default":
		return libp2p.Security(secio.ID, secio.New), nil

	// TODO: add others?
	case "none":
		return libp2p.NoSecurity, nil

	default:
		logger.Errorf("unknown security option: %s", opt)
		return nil, ErrUnknownSecurityOption
	}
}

func bootstrapPeers(ctx context.Context, host host.Host, peers []string) error {
	lenPeers := len(peers)

	for idx, p := range peers {
		logger.Infof("Connecting to %d of %d peers: %s", idx+1, lenPeers, p)
		targetAddr, err := multiaddr.NewMultiaddr(p)
		if err != nil {
			logger.Errorf("err parsing targetAddr from multiaddr\n%v", err)
			return err
		}

		targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
		if err != nil {
			logger.Errorf("err parsing targetInfo from peer addr\n%v", err)
			return err
		}

		if err = host.Connect(ctx, *targetInfo); err != nil {
			logger.Errorf("err connecting\n%v", err)
			return err
		}

		logger.Infof("Connected to %v", targetInfo.ID)
	}

	return nil
}
