package forwards

import (
	"net"
	"sps/responses"
	"sps/util"
)

func HTTPS(client *net.TCPConn, informations []string) {
	if MatchFilter(informations[1]) {
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
