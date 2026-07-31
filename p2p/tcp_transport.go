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

type TCPTransport struct {
	listenAddress string

	listener net.Listener

	HandshakeFunc HandshakeFunc

	decoder Decoder

	mutex sync.RWMutex

	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{

		HandshakeFunc: NOPHandshakeFunc,
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {

	var err error

	t.listener, err = net.Listen("tcp", t.listenAddress)

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

type Temp struct{}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {

		conn.Close()
		fmt.Printf("TCP Handshake error : %s\n", err)
		return

	}

	//Read Loop
	msg := &Temp{}
	for {

		if err := t.decoder.Decode(conn, msg); err != nil {

			fmt.Printf("TCP Error : %s\n", err)
			continue

		}

	}

}
