package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Port       = 8888
	BufferSize = 1 << 16 // 32 KiB
)

var ConnectionClose = "\r\nConnection: close\r\n\r\n"

var RedirectToLocalhost = []byte("HTTP/1.1 301 Moved Permanently\r\nLocation: http://localhost/" + ConnectionClose)
var MethodNotImplemented = []byte("HTTP/1.1 501 Not Implemented" + ConnectionClose)
var Unavailable = []byte("HTTP/1.1 503 Service Unavailable" + ConnectionClose)

var filter = map[string]bool{}

func SimpleForward(from, to net.Conn, isClosed chan bool) {
	fmt.Printf(
		"Forwarding data from %v to %v...\n",
		from.RemoteAddr(),
		to.RemoteAddr(),
	)
	defer from.Close()
	defer to.Close()
	defer fmt.Printf(
		"Closed forward from %v to %v!\n",
		from.RemoteAddr(),
		to.RemoteAddr(),
	)
	sleepTime := 5 * time.Millisecond // 10MB/s
	for {
		data, err := AsyncReceiver(from, BufferSize)
		select {
		case <-isClosed:
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

func ProccessRequest(client *net.TCPConn) {
	fmt.Printf("Receiving header... ")
	header, err := ReadLineFromConnection(client)
	if err != nil {
		fmt.Println("Failed!")
		client.Close()
		return
	}
	fmt.Println("Done")
	informations := strings.Split(header, " ")
	method := informations[0]
	if method == "CONNECT" {
		fmt.Printf("Method %s not implemented!\n", method)
		client.Write(MethodNotImplemented)
		return
	}
	urlParts := strings.Split(informations[1], "://")
	domain := urlParts[0]
	if len(urlParts) == 2 {
		domain = strings.Split(urlParts[1], "/")[0]
	}
	path := strings.Join(strings.Split(urlParts[1], "/")[1:], "/")
	newHeader := fmt.Sprintf("%s /%s %s\r\n", method, path, informations[2])
	port := 80
	domainParts := strings.Split(domain, ":")
	if len(domainParts) == 2 {
		domain = domainParts[0]
		i, _ := strconv.Atoi(domainParts[1])
		if i > 0 {
			port = i
		}
	}
	if _, ok := filter[domain]; ok {
		fmt.Println("URL found on filter!")
		client.Write(RedirectToLocalhost)
		return
	}
	server, err := net.Dial("tcp", fmt.Sprintf("%s:%d", domain, port))
	if err != nil {
		fmt.Println("Failed to connect to remote server!")
		client.Write(Unavailable)
		return
	}
	server.Write([]byte(newHeader))
	isClosed := make(chan bool)
	go SimpleForward(client, server, isClosed)
	go SimpleForward(server, client, isClosed)
}

func ParseFilterFile(f string) {
	data, err := os.ReadFile(f)
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
