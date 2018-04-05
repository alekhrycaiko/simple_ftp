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
	path    []string
	login   bool
	pasv    net.Conn
	ftype   rune
	mode    rune
	input   []string
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

func sendPasv(client *client, address string) {
	msg := []byte(fmt.Sprintf("227 %s\r\n", address))
	client.writer.Write(msg)
	client.writer.Flush()
}
