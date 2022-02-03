package proxy

import (
	"fmt"
	"net"
	"net/url"
	"sps/responses"
	"sps/util"
	"sps/filter"
)

func HTTP(client *net.TCPConn, informations []string, useRegex bool) {
	urlInfo, _ := url.Parse(informations[1])
	port := urlInfo.Port()
	if port == "" {
		port = "80"
	}
	if filter.MatchFilter(informations[1]) {
		Close(client, responses.Filtered)
		return
	}
	server, err := net.Dial("tcp", fmt.Sprintf("%s:%s", urlInfo.Host, port))
	if err != nil {
		Close(client, responses.Unavailable)
		return
	}
	newHeader := fmt.Sprintf(
		"%s %s %s\r\n",
		informations[0],
		urlInfo.Path,
		informations[2],
	)
	server.Write([]byte(newHeader))
	util.Link(client, server)
}
