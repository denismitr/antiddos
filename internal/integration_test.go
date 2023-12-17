package internal

import (
	"context"
	"errors"
	"github.com/denismitr/antiddos/internal/bootstrap"
	"github.com/denismitr/antiddos/internal/quotes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	serverCtx, cancel := context.WithCancel(context.Background())
	s, err := bootstrap.TcpServer(serverCtx, 30, 3, "127.0.0.1", 3333)
	if err != nil {
		t.Fatal(err)
	}

	defer cancel()

	go func() {
		if err := s.Run(serverCtx); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
		}
	}()

	t.Run("client with valid interaction", func(t *testing.T) {
		c := bootstrap.TcpClient(3, 30, "127.0.0.1", 3333)
		conn, closer, err := c.Connect()
		if err != nil {
			t.Fatal(err)
		}
		defer closer()

		clientCtx, clientCancel := context.WithTimeout(serverCtx, 3*time.Second)
		defer clientCancel()

		quote, err := c.Communicate(clientCtx, conn)
		if err != nil {
			t.Fatal(err)
		}

		match := false
		for i := range quotes.Quotes {
			if quotes.Quotes[i] == quote {
				match = true
			}
		}

		assert.Truef(t, match, "wrong quote: [%s]", quote)
	})

	t.Run("client with invalid zeroes", func(t *testing.T) {
		c := bootstrap.TcpClient(2, 30, "127.0.0.1", 3333)
		conn, closer, err := c.Connect()
		if err != nil {
			t.Fatal(err)
		}
		defer closer()

		clientCtx, clientCancel := context.WithTimeout(serverCtx, 3*time.Second)
		defer clientCancel()

		quote, err := c.Communicate(clientCtx, conn)
		if err == nil {
			t.Fatalf("expected an error of non matching zeroes")
		}
		assert.Equal(t, "", quote)
	})
}
