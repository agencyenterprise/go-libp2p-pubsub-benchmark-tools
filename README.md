> libp2p/go-libp2p-pubsub gossip host

# gossip-host
[![Build Status](https://travis-ci.org/agencyenterprise/gossip-host.svg?branch=develop)](https://travis-ci.org/agencyenterprise/gossip-host) [![Coverage Status](https://coveralls.io/repos/github/agencyenterprise/gossip-host/badge.svg?branch=develop)](https://coveralls.io/github/agencyenterprise/gossip-host?branch=develop) [![Go Report Card](https://goreportcard.com/badge/github.com/agencyenterprise/gossip-host)](https://goreportcard.com/report/github.com/agencyenterprise/gossip-host) [![GoDoc](https://godoc.org/github.com/agencyenterprise/gossip-host?status.svg)](https://godoc.org/github.com/agencyenterprise/gossip-host)

TODO


## Usage

1. Check the options in the `configs/host/config.json` file
2. `$ go run ./cmd/host/main.go`
3. Connect a second host to the first: `$ go run ./cmd/host/main.go -l /ip4/127.0.0.2/tcp/3002,/ip4/127.0.0.2/tcp/3003/ws -r :8081 -p <prev. host listen addrs>`. Note, the `-l` flag are the listen addresses. Notice how we've incremented standard local host from `127.0.0.1` to `127.0.0.2`. We could have also simply changed the port address. Also, `-r` is the rpc listen address and needs to be different for this host than the default `:8080`.
4. Send a message to the first host on the rpc channel. The host will then gossip the message to its peers: `$ go run ./cmd/client/main.go -p :8080`


## Commands

### Client

TODO

### Host

TODO

### Subnet

TODO


## Protobufs

This implementation uses RPC to connect to hosts. In order to make the protobuf, you will need to first follow the installation instructions, [here](https://github.com/golang/protobuf), and then run `$ make -C ./pkg/pb/`.


## License

[**MIT**](LICENSE).

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fagencyenterprise%2Fgossip-host.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fagencyenterprise%2Fgossip-host?ref=badge_large)

```
MIT License

Copyright (c) 2019 AE Studio

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
