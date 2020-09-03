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

func TestGroupAndChannel(t *testing.T) {
	ulxd := ULXDReceiver{
		pool: &connpool.Pool{
			TTL:   30 * time.Second,
			Delay: 200 * time.Millisecond,
			NewConnection: func(ctx context.Context) (net.Conn, error) {
				dial := net.Dialer{}
				return dial.DialContext(ctx, "tcp", "TNRB-W121-RCV1.byu.edu:2202")
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	group, channel, err := ulxd.GroupAndChannel(ctx, 1)
	require.NoError(t, err)
	fmt.Printf("group: %d, channel: %d\n", group, channel)
}

func TestTransmitterRFPower(t *testing.T) {
	ulxd := ULXDReceiver{
		pool: &connpool.Pool{
			TTL:   30 * time.Second,
			Delay: 200 * time.Millisecond,
			NewConnection: func(ctx context.Context) (net.Conn, error) {
				dial := net.Dialer{}
				return dial.DialContext(ctx, "tcp", "TNRB-W121-RCV1.byu.edu:2202")
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	power, err := ulxd.TransmitterRFPower(ctx, 1)
	require.NoError(t, err)
	fmt.Printf("transmitter RF power: %v\n", power)
}
