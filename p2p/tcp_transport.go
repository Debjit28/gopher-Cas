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

// this function implements the peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
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
	rpc := RPC{}
	for {

		if err := t.Decoder.Decode(conn, &rpc); err != nil {

			fmt.Printf("TCP Error : %s\n", err)
			return //<---------------earlier had put continue instead of return which caused a mess known an infinite loop of go routines------->

		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc

	}

}
