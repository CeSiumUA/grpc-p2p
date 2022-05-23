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
	addressToCallFrom := conn.LocalAddr().String()
	addressToListenString := ":" + strings.Split(conn.LocalAddr().String(), ":")[1]
	log.Println("closing node connection")
	err = conn.Close()
	if err != nil {
		log.Println("error closing connection")
	}

	log.Println("starting listener:", addressToListenString)
	listener, err := reuseport.Listen("tcp", addressToListenString)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("started listener on", listener.Addr().String())
	defer func(lstnr *net.Listener) {
		if listener != nil {
			log.Println("closing listener")
			err = listener.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}
	}(&listener)
	go func() {
		for {
			log.Println("waiting for clients to accept...")
			connectedPeer, err := listener.Accept()
			if err != nil {
				log.Println("error getting new connection", err)
				continue
			}
			log.Println("new client peer connected", connectedPeer.RemoteAddr().String())
		}

	}()
	log.Println("trying to dial to", addresses[0], "from:", addressToCallFrom)
	peerConnection, err := reuseport.Dial("tcp", addressToCallFrom, addresses[0])
	if err != nil {
		log.Println("error connecting to peer", err)
		return
	}
	fmt.Println("connected to peer successfully", peerConnection.RemoteAddr().String())
}
