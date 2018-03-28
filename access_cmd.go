package main

import (
	"fmt"
	"math/rand"
	"net"
)

func handleAccessCommandUser(conn net.Conn, arg string, loggedIn bool) {
	if arg == "anonymous" && !loggedIn {
		conn.Write([]byte("331 Please specify the password.\r\n"))
	} else if loggedIn {
		conn.Write([]byte("331 Cannot change from anonymous user"))
	} else {
		conn.Write([]byte("530 Anonymous only.\r\n"))
	}
}

func handleAccessCommandPass(conn net.Conn, loggedIn *bool) {
	if !*loggedIn {
		conn.Write([]byte("230 Successful login.\r\n"))
		*loggedIn = true
	} else {
		conn.Write([]byte("230 Already logged in.\r\n"))
	}
}

func handleAccessCommandQuit(conn net.Conn) {
	conn.Write(formatMsg(221, "Quit Successful"))
	conn.Close()
}

func handlePasvConnection(conn net.Conn, pasvConn net.Conn) {
	conn.Write(formatMsg(100, "Doing stuff "))
	// TODO: Handle what happens inside a PASV connection.
}

func setupPasvConnection(conn net.Conn) {
	port := rand.Int()%49151 + 1024
	address := fmt.Sprintf("0.0.0.0:%d", port)
	lstraddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		conn.Write(formatMsg(500, "Invalid response"))
		return
	}
	listener, err := net.ListenTCP("tcp", lstraddr)
	msg := fmt.Sprintf("Connect to %s", address)
	conn.Write(formatMsg(227, msg))
	// TODO extended pasv:
	for {
		pasvConn, err := listener.Accept()
		if err != nil {
			defer conn.Close()
		}
		go handlePasvConnection(conn, pasvConn)
	}
}
