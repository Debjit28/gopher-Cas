package main

import (
	"fmt"
	"log"

	"github.com/Debjit28/gopher-Cas/p2p"
)

func main() {
	opts := p2p.TCPTransportOpt{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tr := p2p.NewTCPTransport(opts)

	go func() {
		for {

			msg := <-tr.Consume()
			fmt.Printf("[%s]: %s\n", msg.From, string(msg.Payload)) //----------------------cleaner and readable output------------
		}

	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	// Keep main process alive while accepting connections
	select {}

}
