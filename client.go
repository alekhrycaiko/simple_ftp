package main

import (
	"bufio"
	"fmt"
	"net"
)

type client struct {
	conn    net.Conn
	writer  *bufio.Writer
	scanner *bufio.Scanner
	dir     string
	login   bool
	pasv    net.Conn
}

func sendMessage(client *client, code int) error {
	msg := []byte(fmt.Sprintf("%d %s\r\n", code, codeMap[code]))
	_, err := client.writer.Write(msg)
	// TODO: best way of handling these errors?
	if err != nil {
		return err
	}
	err = client.writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

// TODO: re-factor this.
func sendPasv(client *client, address string) {
	msg := []byte(fmt.Sprintf("227 %s\r\n", address))
	client.writer.Write(msg)
	client.writer.Flush()
}
