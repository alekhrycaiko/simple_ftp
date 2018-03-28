package main

import (
	"fmt"
	"math/rand"
	"net"
)

func handleAccessCommandUser(c *client) {
	arg, _ := parseConnInput(c.scanner.Text(), 1)
	fmt.Sprintf("%s\r\n", arg)
	if arg == "anonymous" {
		sendMessage(c, 331)
	} else {
		sendMessage(c, 530)
	}
}

func handleAccessCommandPass(c *client) {
	sendMessage(c, 230)
	c.login = true
}

func handleAccessCommandQuit(c *client) {
	sendMessage(c, 221)
	c.conn.Close()
}

func handlePasvConnection(c *client) {
	// TODO: Handle what happens inside a PASV connection.
}

func setupPasvConnection(c *client) {
	port := rand.Int()%49151 + 1024
	address := fmt.Sprintf("0.0.0.0:%d", port)
	lstraddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		sendMessage(c, 500)
		return
	}
	listener, err := net.ListenTCP("tcp", lstraddr)
	msg := fmt.Sprintf("Connect to %s", address)
	sendPasv(c, msg)
	// TODO extended pasv:
	for {
		pasvConn, err := listener.Accept()
		if err != nil {
			defer c.conn.Close()
		}
		c.pasv = pasvConn
		go handlePasvConnection(c)
	}
}
