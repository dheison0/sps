package util

import "net"

func ReadLinesFromBytes(d []byte) [][]byte {
	result := [][]byte{}
	line := []byte{}
	for _, b := range d {
		if b == '\n' {
			result = append(result, line)
			line = []byte{}
			continue
		} else if b == '\r' {
			continue
		}
		line = append(line, b)
	}
	if len(line) > 0 {
		result = append(result, line)
	}
	return result
}

func ReadLineFromConnection(c net.Conn) (string, error) {
	line := []byte{}
	for {
		b := make([]byte, 1)
		size, err := c.Read(b)
		if err != nil {
			return "", err
		} else if b[0] == '\r' || size == 0 {
			continue
		} else if b[0] == '\n' {
			break
		}
		line = append(line, b[0])
	}
	return string(line), nil
}
