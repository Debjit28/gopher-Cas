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

	mutex sync.RWMutex

	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpt) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpt: opts,
	}
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

type Temp struct{}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {

		conn.Close()
		fmt.Printf("TCP Handshake error : %s\n", err)
		return

	}

	//Read Loop

	//buf := make([]byte, 2000)
	msg := &Message{}
	for {

		// n, err := conn.Read(buf)

		// if err != nil {

		// 	fmt.Printf("TCP Error : %s\n", err)

		// }

		if err := t.Decoder.Decode(conn, msg); err != nil {

			fmt.Printf("TCP Error : %s\n", err)
			continue

		}

		msg.From = conn.RemoteAddr()

		fmt.Printf("message %+v\n", msg)

	}

}
