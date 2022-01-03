package pkg

import (
	"net"
	"sps/pkg/forwards"
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
		forwards.HTTPS(client, informations)
	} else {
		forwards.HTTP(client, informations, useRegex)
	}
}
