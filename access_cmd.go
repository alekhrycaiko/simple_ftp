package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"path"
)

func handleAccessCommandUser(c *client) {
	arg, _ := parseConnInput(c.scanner.Text(), 1)
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
	sendMessage(c, 200)
}

func setupPasvConnection(c *client) {
	port := rand.Int()%49151 + 1024
	p1 := port & 0xff
	p2 := (port >> 8) & 0xff
	// TODO: Change the ip to be dynamic.
	outAddr := fmt.Sprintf("0, 0, 0, 0, %d, %d", p1, p2)
	address := fmt.Sprintf("0.0.0.0:%d", port)
	lstraddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		sendMessage(c, 500)
		return
	}
	listener, err := net.ListenTCP("tcp", lstraddr)
	msg := fmt.Sprintf("Connect to %s", outAddr)
	sendPasv(c, msg)
	for {
		pasvConn, err := listener.Accept()
		if err != nil {
			defer c.conn.Close()
		}
		c.pasv = pasvConn
		go handlePasvConnection(c)
	}
}

/**
* Changes the working directory.
 */
func handleAccessCommandCwd(c *client) {
	// check if cpath is valid.
	arg, err := parseConnInput(c.scanner.Text(), 1)
	if arg == "./" || arg == "../" {
		sendMessage(c, 550)
		return
	}
	if err != nil {
		sendMessage(c, 500)
		return
	}
	if arg == ".." {
		handleAccessCommandCdup(c)
		return
	}
	tempPath := append(c.path, arg)
	pathName := path.Join(tempPath...)
	fmt.Println(pathName)
	fi, err := os.Stat(pathName)
	if err != nil {
		sendMessage(c, 550)
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		c.path = tempPath
		sendMessage(c, 250)
	default:
		sendMessage(c, 550)
	}
	return
}

/**
* Moves path up. But not beyond the parent.
 */
func handleAccessCommandCdup(c *client) {
	// Set to 1 to ensure we always have a root folder as parent present in client path obj.
	if len(c.path) == 1 {
		sendMessage(c, 550)
	} else {
		c.path = c.path[:len(c.path)-1]
		fmt.Println(c.path)
		sendMessage(c, 250)
	}
}
