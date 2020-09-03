package shure

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GroupAndChannel returns the group and channel
func (u *ULXDReceiver) GroupAndChannel(ctx context.Context, channel int) (int, int, error) {
	cmd := []byte(fmt.Sprintf("< GET %d GROUP_CHAN >", channel))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := u.sendCommand(ctx, cmd)
	if err != nil {
		return 0, 0, err
	}

	str := string(resp)
	str = strings.TrimPrefix(str, fmt.Sprintf("< REP %d GROUP_CHAN", channel))
	str = strings.TrimSuffix(str, ">")
	str = strings.TrimSpace(str)

	if str == "--,--" {
		return 0, 0, fmt.Errorf("receiver is on a frequency that does not line up with a group and channel")
	}

	split := strings.Split(str, ",")
	if len(split) != 2 {
		return 0, 0, fmt.Errorf("unexpected response format: %q", str)
	}

	group, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, 0, fmt.Errorf("unable to parse group: %w", err)
	}

	ch, err := strconv.Atoi(split[1])
	if err != nil {
		return 0, 0, fmt.Errorf("unable to parse channel: %w", err)
	}

	return group, ch, nil
}

// TransmitterRFPower returns the transmitter RF power
func (u *ULXDReceiver) TransmitterRFPower(ctx context.Context, channel int) (string, error) {
	cmd := []byte(fmt.Sprintf("< GET %d TX_RF_PWR >", channel))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := u.sendCommand(ctx, cmd)
	if err != nil {
		return "", err
	}

	str := string(resp)
	str = strings.TrimPrefix(str, fmt.Sprintf("< REP %d TX_RF_PWR", channel))
	str = strings.TrimSuffix(str, ">")
	str = strings.TrimSpace(str)

	return str, nil
}
