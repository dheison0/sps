package main

var ConnectionClose = "\r\nConnection: close\r\n\r\n"

var RedirectToLocalhost = []byte("HTTP/1.1 301 Moved Permanently\r\nLocation: http://localhost/" + ConnectionClose)
var Unavailable = []byte("HTTP/1.1 503 Service Unavailable" + ConnectionClose)
var Connected = []byte("HTTP/1.1 200 Connection established\r\n\r\n")
