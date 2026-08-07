package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {

	//conn is the underlying connection of the peer
	conn net.Conn

	//if we dial and retrive a connection  ==> outbound == true
	//if we accept and retrive a connection ==> outbound == false

	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpt struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc

	Decoder Decoder
}

type TCPTransport struct {
	TCPTransportOpt

	listener net.Listener

	rpcch chan RPC

	mutex sync.RWMutex

	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpt) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpt: opts,
		rpcch:           make(chan RPC),
	}
}

// Consume implemented the Transport interface , which will return read-only channel.
// for reading the incoming message received from another peer in the network.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {

	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)

	if err != nil {

		return err

	}

	go t.startAcceptLoop()

	return nil

}

func (t *TCPTransport) startAcceptLoop() {
	for {

		conn, err := t.listener.Accept()

		if err != nil {

			fmt.Printf("TCP accept error : %s\n", err)
		}

		fmt.Printf("new incoming connection %+v\n", conn)

		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {

		conn.Close()
		fmt.Printf("TCP Handshake error : %s\n", err)
		return

	}

	//Read Loop

	//buf := make([]byte, 2000)
	rpc := RPC{}
	for {

		// n, err := conn.Read(buf)

		// if err != nil {

		// 	fmt.Printf("TCP Error : %s\n", err)

		// }

		if err := t.Decoder.Decode(conn, &rpc); err != nil {

			fmt.Printf("TCP Error : %s\n", err)
			continue

		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc

	}

}
