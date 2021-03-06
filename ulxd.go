/*
Package shure provides control of Shure ULXD digital wireless systems.
Documentation of Shure's API can be found at http://www.shure.com/americas/support/find-an-answer/ulx-d-crestron-amx-control-strings.

Supported Devices

This is a list of devices that BYU currently uses in production, monitored with this driver.

  Shure ULXD4 Receiver https://www.shure.com/en-US/products/wireless-systems/ulx-d_digital_wireless/ulxd4
  Shure ULXD1 Transmitter https://www.shure.com/en-US/products/wireless-systems/ulx-d_digital_wireless/ulxd1

This list is not comprehensive.
*/
package shure

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/connpool"
)

// ULXDReceiver represents a ULXD Receiver.
type ULXDReceiver struct {
	address string
	pool    *connpool.Pool
}

func NewReceiver(address string) (*ULXDReceiver, error) {
	p := connpool.Pool{
		TTL:   30 * time.Second,
		Delay: 200 * time.Millisecond,
	}

	p.NewConnection = func(ctx context.Context) (net.Conn, error) {
		timeout := 30 * time.Second // default timeout
		// Set timeout according to context if a deadline exists
		if d, ok := ctx.Deadline(); ok {
			timeout = time.Until(d)
		}
		return net.DialTimeout("tcp", fmt.Sprintf("%s:2202", address), timeout)
	}

	return &ULXDReceiver{
		address: address,
		pool:    &p,
	}, nil
}

func (u *ULXDReceiver) sendCommand(ctx context.Context, cmd []byte) ([]byte, error) {
	var resp []byte

	err := u.pool.Do(ctx, func(conn connpool.Conn) error {
		deadline, _ := ctx.Deadline()

		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return fmt.Errorf("unable to write command: %w", err)
		case n != len(cmd):
			return fmt.Errorf("unable to write command: wrote %v/%v bytes", n, len(cmd))
		}

		resp, err = conn.ReadUntil('>', deadline)
		if err != nil {
			return fmt.Errorf("unable to read response: %w", err)
		}

		return nil
	})
	if err != nil {
		return []byte{}, err
	}

	return resp, nil
}
