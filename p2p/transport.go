package p2p

// peer is just that representing the remote node

type Peer interface {
	close() error
}

// Transport is anything that handles communication
// between  the nodes in the network .
// it follow protocols like tcp udp websockets.
type Transport interface {
	ListenAndAccept() error

	Consume() <-chan RPC
}
