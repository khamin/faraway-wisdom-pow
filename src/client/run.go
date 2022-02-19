package client

import (
	"encoding/binary"
	"fmt"
	"net"
)

var BytesOrder = binary.LittleEndian

var read = binary.Read
var write = binary.Write

func Run(server string) ([]byte, error) {
	conn, err := net.Dial("tcp", server)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	ok, err := challenge(conn)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("challenge error")
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}
