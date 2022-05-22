package main

import (
	"fmt"
	"github.com/libp2p/go-reuseport"
	"log"
	"net"
	"strings"
)

func StartP2P() {
	conn, err := net.Dial("tcp4", "5.189.145.4:16574")
	defer func(conn *net.Conn) {
		log.Println("closing connection...")
		if conn != nil {
			err := (*conn).Close()
			if err != nil {
				HandleError(err)
			}
		}
	}(&conn)
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
		addressToListen := ":" + strings.Split(conn.LocalAddr().String(), ":")[1]
		err = conn.Close()
		if err != nil {
			log.Println("error closing connection")
		}
		listener, err := reuseport.Listen("tcp4", addressToListen)
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
		HandleError(err)
		fmt.Println("started listener on", listener.Addr().String())
		peerConnection, err := reuseport.Dial("tcp", listener.Addr().String(), addresses[0])
		if err != nil {
			log.Println("error connecting to peer", err)
			return
		}
		fmt.Println("connected to peer successfully", peerConnection.RemoteAddr().String())
		break
	}
}
