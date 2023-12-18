package protocol

import (
	"encoding/binary"
)

const Delimiter = '#'

type Action uint16

const (
	Request Action = iota
	Challenge
	Solve
	Reject
	Transmit
)

type Payload struct {
	Action Action
	Data   []byte
}

func (p *Payload) Encode() ([]byte, error) {
	buf := make([]byte, 2+2+len(p.Data)+1)
	binary.LittleEndian.PutUint16(buf, uint16(p.Action))
	// todo: verify that data length is not above uint16
	binary.LittleEndian.PutUint16(buf[2:], uint16(len(p.Data)))
	copy(buf[4:], p.Data)
	buf[len(buf)-1] = Delimiter
	return buf, nil
}

func Decode(b []byte) (*Payload, error) {
	p := Payload{}
	p.Action = Action(binary.LittleEndian.Uint16(b))
	length := binary.LittleEndian.Uint16(b[2:])
	p.Data = make([]byte, length)
	copy(p.Data, b[4:])
	return &p, nil
}
