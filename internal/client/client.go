package client

import (
	"bufio"
	"context"
	"fmt"
	"github.com/denismitr/antiddos/internal/protocol"
	"log/slog"
	"net"
	"time"
)

type solver interface {
	Solve(header string) (string, error)
}

type Client struct {
	addr string
	s    solver
}

func New(addr string, s solver) *Client {
	return &Client{
		addr: addr,
		s:    s,
	}
}

func (c *Client) Run(ctx context.Context) error {
	conn, closer, err := c.Connect()
	if err != nil {
		return err
	}
	defer closer()

	slog.Info("client connected to", "addr", c.addr)

	t := time.NewTicker(3 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			if quote, err := c.Communicate(ctx, conn); err != nil {
				return err
			} else {
				slog.With("quote", quote).Info("server transmitted")
			}
		}
	}
}

func (c *Client) Connect() (net.Conn, func() error, error) {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial %s: %w", c.addr, err)
	}

	return conn, conn.Close, nil
}

func (c *Client) Communicate(ctx context.Context, conn net.Conn) (string, error) {
	r := bufio.NewReader(conn)

	if err := c.askForChallenge(ctx, conn); err != nil {
		return "", err
	}

	header, err := c.receiveChallenge(ctx, r)
	if err != nil {
		return "", err
	}

	solution, err := c.doProofOfWork(ctx, header)
	if err != nil {
		return "", err
	}

	if err := c.sendSolution(ctx, solution, conn); err != nil {
		return "", err
	}

	resp, err := c.readTransmission(ctx, r)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (c *Client) askForChallenge(ctx context.Context, conn net.Conn) error {
	slog.Info("asking for a challenge")

	p := protocol.Payload{
		Action: protocol.Request,
	}
	b, err := p.Encode()
	if err != nil {
		return fmt.Errorf("client.askForChallenge payload encode failed: %w", err)
	}
	if _, err := conn.Write(b); err != nil {
		return fmt.Errorf("client.askForChallenge conn.Write failed: %w", err)
	}
	return nil
}

func (c *Client) doProofOfWork(ctx context.Context, header string) (string, error) {
	slog.With("header", header).Info("doing the proof of work on")
	h, err := c.s.Solve(header)
	if err != nil {
		return "", fmt.Errorf("proof of work failed: %w", err)
	}
	return h, nil
}

func (c *Client) sendSolution(_ context.Context, solution string, conn net.Conn) error {
	p := protocol.Payload{
		Action: protocol.Solve,
		Data:   []byte(solution),
	}

	if err := protocol.Send(&p, conn); err != nil {
		return fmt.Errorf("client.Client.sendSolution failed: %w", err)
	}

	return nil
}

func (c *Client) readTransmission(ctx context.Context, r *bufio.Reader) (string, error) {
	resp, err := r.ReadBytes('#')
	if err != nil {
		return "", fmt.Errorf("client.Client.readQoute failed to read bytes: %w", err)
	}

	p, err := protocol.Decode(resp)
	if err != nil {
		return "", fmt.Errorf("client.Client.readQoute failed to decode payload: %w", err)
	}

	switch p.Action {
	case protocol.Reject:
		return "", fmt.Errorf("server rejected the solution")
	case protocol.Transmit:
		return string(p.Data), nil
	default:
		return "", fmt.Errorf("client.Client.readTransmission received unexpected [%d] action", p.Action)
	}
}

func (c *Client) receiveChallenge(ctx context.Context, r *bufio.Reader) (string, error) {
	resp, err := r.ReadBytes('#')
	if err != nil {
		return "", fmt.Errorf("client.askForChallenge read challange resp failed: %w", err)
	}

	slog.Info("challenge received")
	respPayload, err := protocol.Decode(resp)
	if err != nil {
		return "", fmt.Errorf("client.askForChallenge decode resp failed: %w", err)
	}

	if respPayload.Action != protocol.Challenge {
		return "", fmt.Errorf("client.askForChallenge invalid resp payload action %v", respPayload.Action)
	}

	slog.Info("inspecting", "challenge", string(respPayload.Data))
	return string(respPayload.Data), nil
}
