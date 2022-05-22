package main

import (
	"io"
	"log"
	"net"
	"strings"
)

var clientsConnections []*net.Conn = make([]*net.Conn, 0)

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

		go func() {
			defer func() {
				err := conn.Close()
				if err != nil {
					log.Println("error closing connection", err)
					return
				}
			}()
			if err != nil {
				log.Println("error accepting connection", err)
				return
			}
			log.Println("remote client connected", conn.RemoteAddr().String())
			clientsConnections = append(clientsConnections, &conn)
			notifyAllClients()
			oneByte := make([]byte, 1)
			readBytesCount, readErr := conn.Read(oneByte)
			if readErr == io.EOF {
				log.Println("client", conn.RemoteAddr(), "disconnected")
			}
			log.Println("got", readBytesCount, "bytes from connection")
		}()
	}
}

func notifyAllClients() {
	for _, conn := range clientsConnections {
		addr := (*conn).RemoteAddr().String()
		addresses := getAddressesForClient(addr)
		if len(addresses) > 0 {
			addressesStrings := strings.Join(addresses, ",")
			bytes, err := (*conn).Write([]byte(addressesStrings))
			log.Println("wrote", bytes, "bytes to", addr, "data:", addressesStrings)
			if err != nil {
				log.Println("error writing to client")
			}
		}
	}
}

func getAddressesForClient(address string) []string {
	addresses := make([]string, 0)
	for _, v := range clientsConnections {
		stringAddress := (*v).RemoteAddr().String()
		if stringAddress != address {
			addresses = append(addresses, stringAddress)
		}
	}
	return addresses
}
