package proxy

import (
	"net"
	"sps/util"
	"strings"
)

/// ProccessRequest try to choose between forward HTTP or HTTPS
func ProccessRequest(client *net.TCPConn, useRegex bool) {
	header, err := util.ReadLineFromConnection(client)
	if err != nil {
		client.Close()
		return
	}
	informations := strings.Split(header, " ")
	method := informations[0]
	if method == "CONNECT" {
		HTTPS(client, informations)
	} else {
		HTTP(client, informations, useRegex)
	}
}
