package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dytlzl/tervi/internal/debug"
)

func main() {
	// UDPサーバ
	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP("localhost"),
		Port: debug.UDP_PORT,
	}
	updLn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 1024)
	log.Println("Starting UDP Server...")

	for {
		n, _, err := updLn.ReadFromUDP(buf)
		if err != nil {
			log.Fatalln(err)
		}

		go func() {
			fmt.Printf(string(buf[:n]))
		}()
	}
}
