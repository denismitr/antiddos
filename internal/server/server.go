package server

import (
	"bufio"
	"context"
	"fmt"
	"github.com/denismitr/antiddos/internal/protocol"
	"io"
	"log/slog"
	"net"
)

type requestHandler interface {
	Handle(
		ctx context.Context,
		req []byte,
		clientIP string,
	) (*protocol.Payload, error)
}

type Server struct {
	addr string
	rh   requestHandler
}

func New(addr string, h requestHandler) *Server {
	return &Server{
		addr: addr,
		rh:   h,
	}
}

func (s *Server) Run(ctx context.Context) error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("server failed to start listening on %s: %w", s.addr, err)
	}
	defer l.Close()

	slog.With("tcp", s.addr).Info("listening on address")

	errCh := make(chan error, 1)
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				errCh <- fmt.Errorf("failed to accept a new connection: %w", err)
				return
			}

			go s.handleConnection(ctx, conn)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	slog.With("address", conn.RemoteAddr().String()).Info("new client")
	defer conn.Close()

	r := bufio.NewReader(conn)

	for {
		if ctx.Err() != nil {
			slog.Error(ctx.Err().Error())
			return
		}

		b, err := r.ReadBytes(protocol.Delimiter)
		if err != nil {
			if err == io.EOF {
				slog.Info("connection ended")
				return
			}

			slog.With("error", err.Error()).Error("server.Server.handleConnection failed to read payload")
			return
		}

		payload, err := s.rh.Handle(ctx, b, conn.RemoteAddr().String())
		if err != nil {
			slog.With("error", err.Error()).Error("server.Server.handleConnection failed to process request")
			return
		}

		if err := protocol.Send(payload, conn); err != nil {
			slog.
				With("error", err.Error()).
				With("client address", conn.RemoteAddr().String()).
				Error("server failed to send payload")
		}
	}
}
