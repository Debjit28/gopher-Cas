package p2p

import "net"

// Message represent that is being sent over each transport between 2 nodes in za network
type RPC struct {
	From    net.Addr
	Payload []byte
}
