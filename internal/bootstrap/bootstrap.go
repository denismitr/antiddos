package bootstrap

import (
	"context"
	"fmt"
	"github.com/denismitr/antiddos/internal/challenge"
	"github.com/denismitr/antiddos/internal/client"
	"github.com/denismitr/antiddos/internal/protocol"
	"github.com/denismitr/antiddos/internal/quotes"
	"github.com/denismitr/antiddos/internal/server"
	"github.com/denismitr/antiddos/internal/store/adapters/embedded"
	"github.com/denismitr/antiddos/internal/store/adapters/nope"
)

func TcpServer(
	ctx context.Context,
	maxDuration uint64,
	zeroes uint8,
	host string, port int,
) (*server.Server, error) {
	store, err := embedded.New(ctx, maxDuration)
	if err != nil {
		return nil, err
	}

	c := challenge.New(store, zeroes, maxDuration)
	p := protocol.New(c, quotes.New())
	addr := fmt.Sprintf("%s:%d", host, port)
	return server.New(addr, p), nil
}

func TcpClient(zeroes uint8, maxDuration uint64, host string, port int) *client.Client {
	addr := fmt.Sprintf("%s:%d", host, port)
	clientSideValidator := nope.Nope{}
	solver := challenge.New(clientSideValidator, zeroes, maxDuration)
	return client.New(addr, solver)
}
