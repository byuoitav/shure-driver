package shure

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/byuoitav/connpool"
	"github.com/stretchr/testify/require"
)

func TestTransmitterType(t *testing.T) {
	ulxd := ULXDReceiver{
		pool: &connpool.Pool{
			TTL:   30 * time.Second,
			Delay: 200 * time.Millisecond,
			NewConnection: func(ctx context.Context) (net.Conn, error) {
				dial := net.Dialer{}
				return dial.DialContext(ctx, "tcp", "CB-254-RCV1.byu.edu:2202")
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	typ, err := ulxd.TransmitterType(ctx, 1)
	require.NoError(t, err)
	fmt.Printf("transmitter type: %v\n", typ)
}
