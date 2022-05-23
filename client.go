package main

import (
	"fmt"
	"github.com/libp2p/go-reuseport"
	"log"
	"net"
	"strings"
)

func StartP2P() {
	conn, err := net.Dial("tcp", "5.189.145.4:16574")
	defer func(conn *net.Conn) {
		log.Println("closing connection...")
		if conn != nil {
			err := (*conn).Close()
			if err != nil {
				log.Fatalln(err)
			}
		}
	}(&conn)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("client connected, local:", conn.LocalAddr(), "remote:", conn.RemoteAddr())

	readBytes := make([]byte, 2048)
	readBytesCount, err := conn.Read(readBytes)
	if err != nil {
		fmt.Println("error reading data from connection")
		return
	}
	str := string(readBytes[:readBytesCount])
	addresses := strings.Split(str, ",")
	log.Println("got addresses from server:", addresses)
	addressToListen := "0.0.0.0:" + strings.Split(conn.LocalAddr().String(), ":")[1]
	err = conn.Close()
	if err != nil {
		log.Println("error closing connection")
	}
	log.Println("starting listener:", addressToListen)
	listener, err := reuseport.Listen("tcp", addressToListen)
	if err != nil {
		log.Fatalln(err)
	}
	go func(lstnr *net.Listener) {
		for {
			connectedPeer, err := (*lstnr).Accept()
			if err != nil {
				log.Println("error getting new connection", err)
				continue
			}
			log.Println("new client peer connected", connectedPeer.RemoteAddr().String())
		}
	}(&listener)
	fmt.Println("started listener on", listener.Addr().String())
	peerConnection, err := reuseport.Dial("tcp", listener.Addr().String(), addresses[0])
	if err != nil {
		log.Println("error connecting to peer", err)
		return
	}
	fmt.Println("connected to peer successfully", peerConnection.RemoteAddr().String())
}
