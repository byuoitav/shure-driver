package shure

import (
	"bufio"
	"net"
	"time"
)

type AudioControl struct {
	Address string
}

type Connection struct {
	Conn net.Conn
}

//GetConnection will form a connection with a device at Address
func (s *AudioControl) GetConnection() (*Connection, error) {
	connection, err := net.DialTimeout("tcp", s.Address+":2202", time.Second*3)
	if err != nil {
		return nil, err
	}

	return &Connection{
		Conn: connection,
	}, nil
}

//ReadEvent will read an event from a shure audio device
func (c *Connection) ReadEvent() (string, error) {
	reader := bufio.NewReader(c.Conn)

	data, err := reader.ReadString('>') // does readstring block?
	if err != nil {
		return "", err
	}

	return data, nil
}
