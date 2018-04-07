package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

// handleType provides support for the type command.
func handleType(c *client) {
	if len(c.input) < 2 {
		sendMessage(c, 500)
		return
	}
	arg := rune(c.input[1][0])
	// Only supporting one mode currently.
	if arg == 'A' {
		c.ftype = arg
		sendMessage(c, 200)
	} else {
		sendMessage(c, 504)
	}
}

// getPassiveConn sets up passive connection listener on the given address.
func getPassiveConn(c *client, lstraddr *net.TCPAddr) {
	listener, err := net.ListenTCP("tcp", lstraddr)
	if err != nil {
		sendMessage(c, 500)
	}
	listener.SetDeadline(time.Now().Add(2 * time.Minute))
	defer listener.Close()
	for {
		pasvConn, err := listener.Accept()
		if err == err.(*net.OpError) {
			sendMessage(c, 500)
			break
		} else if err != nil {
			continue
		}
		c.pasv = pasvConn
		break
	}

}

// handlePasv sets up address to being called by the client.
func handlePasv(c *client) {
	port := rand.Int()%49151 + 1024
	p1 := port & 0xff
	p2 := (port >> 8) & 0xff
	outAddr := fmt.Sprintf("(0,0,0,0,%d,%d)", p1, p2)
	fmt.Println(port)
	address := fmt.Sprintf("0.0.0.0:%d", port)
	fmt.Println(address)
	lstraddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		sendMessage(c, 500)
		return
	}
	msg := fmt.Sprintf("Connect to %s", outAddr)
	sendPasv(c, msg)
	go getPassiveConn(c, lstraddr)
}

// handleMode sets transfer mode. Currently always reports 200, as nothing special goes on here!
func handleMode(c *client) {
	sendMessage(c, 200)
}
