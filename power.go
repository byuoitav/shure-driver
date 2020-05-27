package shure

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	battChargeOFF          = 255
	battChargeCalculating  = 254
	battRunTimeOFF         = 65535
	battRunTimeCalculating = 65534
	battBarsOFF            = 255
	battBarsCalculating    = 254
)

//GetBatteryCharge requests the battery charge of the device on channel 'channel'
func (c *Connection) GetBatteryCharge(channel int) (string, error) {
	msg := fmt.Sprintf("< GET %d BATT_CHARGE >", channel)

	status, err := c.getBatteryStatus(msg)
	if err != nil {
		return "", fmt.Errorf("failed to get battery charge: %s", err.Error())
	}

	value, err := strconv.Atoi(status)
	if err != nil {
		return "", fmt.Errorf("could not parse response: %s", err.Error())
	}

	if value == battChargeOFF {
		return "off", nil // the transmitter is off
	} else if value == battChargeCalculating {
		return "calculating", nil
	}

	return strconv.Itoa(value), nil
}

//GetBatteryRunTime requests the battery run time of the device on channel 'channel'
func (c *Connection) GetBatteryRunTime(channel int) (string, error) {
	msg := fmt.Sprintf("< GET %d BATT_RUN_TIME >", channel)

	status, err := c.getBatteryStatus(msg)
	if err != nil {
		return "", fmt.Errorf("failed to get battery run time: %s", err.Error())
	}

	value, err := strconv.Atoi(status)
	if err != nil {
		return "", fmt.Errorf("could not parse response: %s", err.Error())
	}

	if value == battRunTimeOFF {
		return "off", nil // the transmitter is off
	} else if value == battRunTimeCalculating {
		return "calculating", nil
	}

	return strconv.Itoa(value), nil
}

//GetBatteryBars requests the battery bars of the device on channel 'channel'
func (c *Connection) GetBatteryBars(channel int) (string, error) {
	msg := fmt.Sprintf("< GET %d BATT_BARS >", channel)

	status, err := c.getBatteryStatus(msg)
	if err != nil {
		return "", fmt.Errorf("failed to get battery bars: %s", err.Error())
	}

	value, err := strconv.Atoi(status)
	if err != nil {
		return "", fmt.Errorf("could not parse response: %s", err.Error())
	}

	if value == battBarsOFF {
		return "off", nil
	} else if value == battBarsCalculating {
		return "calculating", nil
	}

	return strconv.Itoa(value), nil
}

func (c *Connection) getBatteryStatus(msg string) (string, error) {
	resp, err := c.SendCommand(msg)
	if err != nil {
		return "", fmt.Errorf("failed to get battery status: %s", err.Error())
	}

	re := regexp.MustCompile(`\d+ >`)
	val := re.FindString(resp)
	val = strings.TrimSuffix(val, " >")

	return val, nil
}

//GetPowerStatus requests the power status of the device on channel 'channel'
func (c *Connection) GetPowerStatus(channel int) (string, error) {
	msg := fmt.Sprintf("< GET %d TX_TYPE >", channel)

	resp, err := c.SendCommand(msg)
	if err != nil {
		return "", fmt.Errorf("failed to get power status: %s", err.Error())
	}

	if !strings.Contains(resp, "TX_TYPE") {
		return "", fmt.Errorf("Invalid response, expected type 'TX_TYPE', received: %s", resp)
	} else if strings.Contains(resp, "UNKN") {
		return "standby", nil
	}

	return "on", nil
}
