package shure

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// TransmitterType returns the type of the transmitter
func (u *ULXDReceiver) TransmitterType(ctx context.Context, channel int) (string, error) {
	cmd := []byte(fmt.Sprintf("< GET %d TX_TYPE >", channel))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := u.sendCommand(ctx, cmd)
	if err != nil {
		return "", err
	}

	str := string(resp)
	str = strings.TrimPrefix(str, fmt.Sprintf("< REP %d TX_TYPE", channel))
	str = strings.TrimSuffix(str, ">")
	str = strings.TrimSpace(str)

	return str, nil
}

