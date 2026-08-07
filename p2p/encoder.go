package p2p

import (
	"bytes"
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct{} // <-------basic encoder and decoder in golang ------------------->

func (dec GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg) // <---------------this decodes in raw bytes matlab numbers dega ------------>
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error {

	buf := make([]byte, 2000) //<-----------temporary limit hai ----------------------->

	n, err := r.Read(buf) //<------------------- mssg lenght and err is well error my frnd -------------->

	if n == 0 && err != nil { //<-------------------in case 0 aaara ha ho uske liye hai and error nil na ho return error why makes debugging easier--->

		return err
	}

	msg.Payload = bytes.TrimSpace(buf[:n]) //<---------------------------makes output much more cleaner and easier to understand ----------->

	return nil
}
