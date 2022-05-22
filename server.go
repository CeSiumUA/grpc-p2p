package main

import (
	"log"
	"net"
)

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
	}
}
