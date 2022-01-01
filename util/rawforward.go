package util

import (
	"fmt"
	"net"
	"time"
)

const BufferSize = 1 << 16 // 64 KiB

func AsyncReceiver(c net.Conn, bs uint) (chan []byte, chan error) {
	data := make(chan []byte)
	err := make(chan error)
	go func() {
		received := make([]byte, bs)
		s, e := c.Read(received)
		if e != nil {
			err <- e
		} else {
			data <- received[:s]
		}
	}()
	return data, err
}

/// RawForward forwards raw data between two connections
func RawForward(from, to net.Conn, isClosed chan bool) {
	defer func() {
		from.Close()
		to.Close()
	}()
	sleepTime := 1 * time.Millisecond
	for {
		data, err := AsyncReceiver(from, BufferSize)
		select {
		case <-isClosed:
			fmt.Printf(
				"Closed forward from %v to %v!\n",
				from.RemoteAddr(),
				to.RemoteAddr(),
			)
			return
		case <-err:
			isClosed <- true
			return
		case d := <-data:
			_, e := to.Write(d)
			if e != nil {
				isClosed <- true
				return
			}
		}
		time.Sleep(sleepTime)
	}
}

/// Link links two connections using the SimpleForward function
func Link(client, server net.Conn) {
	fmt.Printf(
		"Forwarding data from %v to %v...\n",
		server.RemoteAddr(),
		client.RemoteAddr(),
	)
	isClosed := make(chan bool)
	go RawForward(client, server, isClosed)
	go RawForward(server, client, isClosed)
}
