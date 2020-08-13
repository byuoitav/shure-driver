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

func TestBatteryCharge(t *testing.T) {
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

	percentage, err := ulxd.BatteryCharge(ctx, 1)
	require.NoError(t, err)
	fmt.Printf("battery charge: %v\n", percentage)
}

func TestBatteryRunTime(t *testing.T) {
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

	left, err := ulxd.BatteryRunTime(ctx, 1)
	require.NoError(t, err)
	fmt.Printf("battery run time: %v\n", left)
}

func TestBatteryType(t *testing.T) {
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

	typ, err := ulxd.BatteryType(ctx, 1)
	require.NoError(t, err)
	fmt.Printf("battery type: %v\n", typ)
}
