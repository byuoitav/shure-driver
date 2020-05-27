package shure

import (
	"bufio"
	"fmt"
	"net"
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

	return data, nil
}

//SendCommand sends a command to the shure receiver and returns the response
func (c *Connection) SendCommand(msg string) (string, error) {
	c.Conn.Write([]byte(msg))

	resp, err := c.ReadEvent()
	if err != nil {
		return "", fmt.Errorf("failed to read response: %s", err.Error())
	}

	return resp, nil
}
