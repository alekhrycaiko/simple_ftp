package main

import (
	"fmt"
	"math/rand"
	"net"
)

func handleTransferCommandType(c *client) {
	if len(c.input) < 2 {
		sendMessage(c, 500)
		return
	}
	arg := rune(c.input[1][0])
	if arg == 'A' || arg == 'I' {
		c.ftype = arg
		sendMessage(c, 200)
	} else {
		sendMessage(c, 504)
	}
}

func handleTransferCommandPasv(c *client) {
	port := rand.Int()%49151 + 1024
	p1 := port & 0xff
	p2 := (port >> 8) & 0xff
	// TODO: Change the ip to be dynamic.
	outAddr := fmt.Sprintf("(0,0,0,0,%d,%d)", p1, p2)
	fmt.Println(port)
	address := fmt.Sprintf("0.0.0.0:%d", port)
	fmt.Println(address)
	lstraddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		sendMessage(c, 500)
		return
	}
	listener, err := net.ListenTCP("tcp", lstraddr)
	if err != nil {
		sendMessage(c, 500)
	}
	msg := fmt.Sprintf("Connect to %s", outAddr)
	sendPasv(c, msg)
	// TODO: set up a timeout.
	for {
		pasvConn, err := listener.Accept()
		if err != nil {
			continue
		}
		c.pasv = pasvConn
		break
	}
	return
}

// Handles what mode we set the transfer to, for now only support stream.
func handleTransferCommandMode(c *client) {
	sendMessage(c, 200)
}
