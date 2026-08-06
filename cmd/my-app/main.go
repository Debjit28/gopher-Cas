package main

import (
	"log"

	"github.com/Debjit28/gopher-Cas/p2p"
)

func main() {
	opts := p2p.TCPTransportOpt{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.GOBDecoder{},
	}

	tr := p2p.NewTCPTransport(opts)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	// Keep main process alive while accepting connections
	select {}

}
