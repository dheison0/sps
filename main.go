package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	Port       = 8888
	BufferSize = 1 << 16 // 64 KiB
)

var filter = map[string]bool{}

/// SimpleForward forwards raw data between two connections
func SimpleForward(from, to net.Conn, isClosed chan bool) {
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
	go SimpleForward(client, server, isClosed)
	go SimpleForward(server, client, isClosed)
}

func ForwardHTTP(client *net.TCPConn, informations []string) {
	urlInfo, _ := url.Parse(informations[1])
	port := urlInfo.Port()
	if port == "" {
		port = "80"
	}
	if _, ok := filter[urlInfo.Host]; ok {
		fmt.Printf("%s blocked!\n", urlInfo.Host)
		client.Write(Filtered)
		return
	}
	server, err := net.Dial("tcp", fmt.Sprintf("%s:%s", urlInfo.Host, port))
	if err != nil {
		client.Write(Unavailable)
		return
	}
	newHeader := fmt.Sprintf(
		"%s %s %s\r\n",
		informations[0],
		urlInfo.Path,
		informations[2],
	)
	server.Write([]byte(newHeader))
	Link(client, server)
}

func ForwardHTTPS(client *net.TCPConn, informations []string) {
	domain := strings.Split(informations[1], ":")[0]
	if _, ok := filter[domain]; ok {
		fmt.Printf("%s blocked!\n", domain)
		client.Write(Filtered)
		return
	}
	server, err := net.Dial("tcp", informations[1])
	if err != nil {
		client.Write(Unavailable)
		return
	}
	// Reads the headers received from the client so as not to forward
	// them to the server
	for {
		line, err := ReadLineFromConnection(client)
		if err != nil {
			client.Close()
			return
		}
		if line == "" {
			break
		}
	}
	client.Write(Connected)
	Link(client, server)
}

/// ProccessRequest try to choose between forward HTTP or HTTPS
func ProccessRequest(client *net.TCPConn) {
	header, err := ReadLineFromConnection(client)
	if err != nil {
		client.Close()
		return
	}
	informations := strings.Split(header, " ")
	method := informations[0]
	if method == "CONNECT" {
		ForwardHTTPS(client, informations)
	} else {
		ForwardHTTP(client, informations)
	}
}

func ParseFilterFile(f string) {
	data, err := ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := ReadLinesFromBytes(data)
	for _, url := range lines {
		filter[string(url)] = true
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("usage: %s filter_file.txt\n", os.Args[0])
		os.Exit(1)
	}
	ParseFilterFile(os.Args[1])
	server, err := net.ListenTCP("tcp", &net.TCPAddr{Port: Port})
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	fmt.Printf("Server started at port %d!\n", Port)
	for {
		client, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go ProccessRequest(client)
	}
}
