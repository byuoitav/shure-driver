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
	Conn net.Conn
}

//GetConnection will form a connection with a device at Address
func (s *AudioControl) GetConnection() (*Connection, error) {
	connection, err := net.Dial("tcp", s.Address+":2202")
	if err != nil {
		return nil, fmt.Errorf("can't make connection to receiver, %s", err.Error())
	}

	return &Connection{
		Conn: connection,
	}, nil
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
