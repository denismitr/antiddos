package protocol

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
)

var (
	ErrInvalidRequestAction = errors.New("invalid request action")
)

type challenger interface {
	Create(string) (string, error)
	Solve(header string) (string, error)
}

type transmissionProvider interface {
	Provide() string
}

type Protocol struct {
	c  challenger
	tp transmissionProvider
}

func New(c challenger, tp transmissionProvider) *Protocol {
	return &Protocol{
		c:  c,
		tp: tp,
	}
}

func (pr *Protocol) Handle(_ context.Context, req []byte, clientIP string) (*Payload, error) {
	p, err := Decode(req)
	if err != nil {
		return nil, fmt.Errorf("failed to decode incoming request with %s: %w", string(req), err)
	}

	switch p.Action {
	case Request:
		d, err := pr.c.Create(clientIP)
		if err != nil {
			return nil, fmt.Errorf("request action failed: %w", err)
		}

		p := Payload{
			Action: Challenge,
			Data:   []byte(d),
		}

		return &p, nil
	case Solve:
		header, err := pr.c.Solve(string(p.Data))
		errWrapped := fmt.Errorf("solve action failed: %w", err)
		if err != nil {
			slog.With("error", errWrapped).Error("rejecting solve")
			return &Payload{
				Action: Reject,
				Data:   []byte(errWrapped.Error()),
			}, nil
		}

		slog.With("header", header).Info("confirmed correct solve")
		transmission := pr.tp.Provide()

		return &Payload{
			Action: Transmit,
			Data:   []byte(transmission),
		}, nil
	default:
		return nil, ErrInvalidRequestAction
	}
}

func Send(p *Payload, w io.Writer) error {
	b, err := p.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("failed to write encoded payload: %w", err)
	}

	return nil
}
