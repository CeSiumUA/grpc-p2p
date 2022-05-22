package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func StartP2P() {
	listener, err := net.Listen("tcp4", ":4566")
	HandleError(err)
	fmt.Println("started listener on", listener.Addr().String())
	laddr, err := net.ResolveTCPAddr("tcp4", ":4566")
	HandleError(err)
	raddr, err := net.ResolveTCPAddr("tcp4", "5.189.145.4:16574")
	HandleError(err)
	conn, err := net.DialTCP("tcp4", laddr, raddr)
	defer func(conn *net.TCPConn) {
		log.Println("closing connection...")
		err := conn.Close()
		if err != nil {
			HandleError(err)
		}
	}(conn)
	HandleError(err)
	log.Println("client connected, local:", conn.LocalAddr(), "remote:", conn.RemoteAddr())
	for {
		readBytes := make([]byte, 2048)
		readBytesCount, err := conn.Read(readBytes)
		if err != nil {
			fmt.Println("error reading data from connection")
			return
		}
		if readBytesCount == 0 {
			continue
		}
		str := string(readBytes[:readBytesCount])
		addresses := strings.Split(str, ",")
		log.Println("got addresses from server:", addresses)
	}
}
