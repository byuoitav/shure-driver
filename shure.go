package shure

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
)

type AudioControl struct {
	Address string
}

type Connection struct {
	Conn    net.Conn
	Address string
}

//GetConnection will form a connection with a device at Address
func (s *AudioControl) GetConnection() (*Connection, error) {
	conn := &Connection{
		Address: s.Address + ":2202",
	}
	err := conn.connect()
	if err != nil {
		return nil, fmt.Errorf("can't make connection to receiver, %s", err.Error())
	}

	return conn, nil
}

func (c *Connection) connect() error {
	connection, err := net.Dial("tcp", c.Address)
	if err != nil {
		return err
	}
	c.Conn = connection
	return nil
}

//ReadEvent will read an event from a shure audio device
func (c *Connection) ReadEvent() (string, error) {
	reader := bufio.NewReader(c.Conn)

	data, err := reader.ReadString('>')
	if err != nil {
		return "", err
	}

	return string(data), nil
}

//GetBatteryCharge requests the battery charge of the device on channel 'channel'
func (c *Connection) GetBatteryCharge(channel int) (string, error) {
	msg := fmt.Sprintf("< GET %d BATT_CHARGE >", channel)

	return c.getBatteryStatus(msg)
}

//GetBatteryRunTime requests the battery run time of the device on channel 'channel'
func (c *Connection) GetBatteryRunTime(channel int) (string, error) {
	msg := fmt.Sprintf("< GET %d BATT_RUN_TIME >", channel)

	return c.getBatteryStatus(msg)
}

//GetBatteryBars requests the battery bars of the device on channel 'channel'
func (c *Connection) GetBatteryBars(channel int) (string, error) {
	msg := fmt.Sprintf("< GET %d BATT_BARS >", channel)

	return c.getBatteryStatus(msg)
}

func (c *Connection) getBatteryStatus(msg string) (string, error) {
	c.Conn.Write([]byte(msg))

	reader := bufio.NewReader(c.Conn)
	resp, err := reader.ReadString('>')
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`\d+ >`)
	val := re.FindString(resp)
	val = strings.TrimSuffix(val, " >")

	return val, nil
}
