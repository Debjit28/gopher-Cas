package main

import (
	"github.com/Debjit28/gopher-Cas/p2p"
)

func main() {
	opts := p2p.TCPTransportOpt{
		ListenAddr: ":3000",
	}

	tr := p2p.NewTCPTransport(opts)

	if err := tr.ListenAndAccept(); err != nil {
		panic(err)
	}

	// Keep main process alive while accepting connections
	select {}

}
