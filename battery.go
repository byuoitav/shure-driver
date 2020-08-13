package shure

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrTransmitterOff = errors.New("transmitter powered off")
	ErrCalculating    = errors.New("calculating")
	ErrOther          = errors.New("unknown error code")
)

// BatteryCharge returns the battery charge as a percentage between [0, 100]
func (u *ULXDReceiver) BatteryCharge(ctx context.Context, channel int) (int, error) {
	cmd := []byte(fmt.Sprintf("< GET %d BATT_CHARGE >", channel))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := u.sendCommand(ctx, cmd)
	if err != nil {
		return 0, err
	}

	str := string(resp)
	str = strings.TrimPrefix(str, fmt.Sprintf("< REP %d BATT_CHARGE", channel))
	str = strings.TrimSuffix(str, ">")

	percentage, err := strconv.Atoi(strings.TrimSpace(str))
	if err != nil {
		return 0, fmt.Errorf("unable to parse response: %w", err)
	}

	switch percentage {
	case 255:
		return 0, ErrTransmitterOff
	case 254:
		return 0, ErrCalculating
	case 253, 252:
		return 0, ErrOther
	}

	return percentage, nil
}

// BatteryRunTime returns the duration until the battery dies
func (u *ULXDReceiver) BatteryRunTime(ctx context.Context, channel int) (time.Duration, error) {
	cmd := []byte(fmt.Sprintf("< GET %d BATT_RUN_TIME >", channel))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := u.sendCommand(ctx, cmd)
	if err != nil {
		return 0, err
	}

	str := string(resp)
	str = strings.TrimPrefix(str, fmt.Sprintf("< REP %d BATT_RUN_TIME", channel))
	str = strings.TrimSuffix(str, ">")

	minutes, err := strconv.Atoi(strings.TrimSpace(str))
	if err != nil {
		return 0, fmt.Errorf("unable to parse response: %w", err)
	}

	switch minutes {
	case 65535:
		return 0, ErrTransmitterOff
	case 65534:
		return 0, ErrCalculating
	case 65533, 65532:
		return 0, ErrOther
	}

	return time.Minute * time.Duration(minutes), nil
}

// BatteryType returns the type of the battery installed
func (u *ULXDReceiver) BatteryType(ctx context.Context, channel int) (string, error) {
	cmd := []byte(fmt.Sprintf("< GET %d BATT_TYPE >", channel))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := u.sendCommand(ctx, cmd)
	if err != nil {
		return "", err
	}

	str := string(resp)
	str = strings.TrimPrefix(str, fmt.Sprintf("< REP %d BATT_TYPE", channel))
	str = strings.TrimSuffix(str, ">")
	str = strings.TrimSpace(str)

	return str, nil
}
