package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
)

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
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(p); err != nil {
		return nil, fmt.Errorf("failed to encode payload: %w", err)
	}
	buf.WriteByte('#')
	return buf.Bytes(), nil
}

func Decode(b []byte) (*Payload, error) {
	buf := bytes.NewReader(b)
	dec := json.NewDecoder(buf)
	p := Payload{}
	if err := dec.Decode(&p); err != nil {
		return nil, fmt.Errorf("could not Decode payload: %w", err)
	}

	return &p, nil
}
