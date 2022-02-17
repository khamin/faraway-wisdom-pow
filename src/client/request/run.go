package request

import "net"

func Run(server string) ([]byte, error) {
	conn, err := net.Dial("tcp", server)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}
