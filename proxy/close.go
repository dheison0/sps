package proxy

import (
	"net"
	"time"
)

func Close(c net.Conn, d []byte) {
	c.Write(d)
	time.Sleep(time.Duration(len(d)) * time.Millisecond)
	c.Close()
}
