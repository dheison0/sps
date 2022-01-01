package forwards

import (
	"fmt"
	"net"
	"sps/responses"
	"sps/util"
	"strings"
)

func HTTPS(client *net.TCPConn, informations []string) {
	domain := strings.Split(informations[1], ":")[0]
	if _, ok := Filter[domain]; ok {
		fmt.Printf("%s blocked!\n", domain)
		client.Write(responses.Filtered)
		return
	}
	server, err := net.Dial("tcp", informations[1])
	if err != nil {
		Close(client, responses.Unavailable)
		return
	}
	// Reads the headers received from the client so as not to forward
	// them to the server
	for {
		line, err := util.ReadLineFromConnection(client)
		if err != nil {
			client.Close()
			return
		}
		if line == "" {
			break
		}
	}
	client.Write(responses.Connected)
	util.Link(client, server)
}
