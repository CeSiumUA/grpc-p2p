package main

import (
	"io"
	"log"
	"net"
)

var clientsConnections []net.Addr = make([]net.Addr, 0)

func Serve() {
	serverListener, err := net.Listen("tcp4", ":16574")
	HandleError(err)
	defer func(serverListener net.Listener) {
		err := serverListener.Close()
		if err != nil {

		}
	}(serverListener)
	for {
		conn, err := serverListener.Accept()
		HandleError(err)
		log.Println("remote client connected", conn.RemoteAddr().String())
		clientsConnections = append(clientsConnections, conn.RemoteAddr())
		oneByte := make([]byte, 1)
		_, readErr := conn.Read(oneByte)
		if readErr == io.EOF {
			log.Println("client disconnected")
		}
	}
}
