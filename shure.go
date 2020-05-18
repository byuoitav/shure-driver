package shure

type AudioControl struct {
}

type Connection struct {
	// similar to net.Conn
}

//GetConnection will form a connection with a shure audio device
func (s *AudioControl) GetConnection() (*Connection, error) {
	return nil, nil
}

//ReadEvent will read an event from a shure audio device
func (c *Connection) ReadEvent() error {
	return nil
}
