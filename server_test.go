package main

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"testing"
	"time"
)

type test_client struct {
	conn   net.Conn
	pasv   net.Conn
	reader *textproto.Reader
	writer *textproto.Writer
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())

}
func connect() *test_client {
	c := &test_client{}
	conn, err := net.DialTimeout("tcp", "0.0.0.0:3000", time.Minute)
	c.conn = conn
	if err != nil {
		fmt.Println("error?")
		fmt.Println(err)
		return nil
	}
	c.reader = textproto.NewReader(bufio.NewReader(c.conn))
	c.writer = textproto.NewWriter(bufio.NewWriter(c.conn))

	c.read(220)
	c.writer.PrintfLine("user anonymous")
	c.read(331)
	c.writer.PrintfLine("pass somepassword")
	c.read(230)
	return c
}
func (c *test_client) read(code int) {
	_, msg, err := c.reader.ReadResponse(code)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error! Reading Code %d failed", code))
	} else {
		fmt.Println(msg)
	}
}

// Tests Cwd and Cdup commands. Assumes we have a file_directory dir available.
func TestCds(m *testing.T) {
	os.Mkdir("file_directory", 0700)
	c := connect()
	c.writer.PrintfLine("cwd file_directory")
	c.read(250)
	c.writer.PrintfLine("cwd ..")
	c.read(250)
	c.writer.PrintfLine("cdup")
	c.read(550)
	c.writer.PrintfLine("cwd file_directory")
	c.read(250)
	c.writer.PrintfLine("cwd file_directory")
	c.read(550)
	c.writer.PrintfLine("cdup")
	c.writer.PrintfLine("quit")
	fmt.Println(c.reader.ReadResponse(0))
	os.Remove("file_directory")
}

// Ensures our FTP server doesnt crash with multi connections.
func TestMultiConnect(b *testing.T) {
	for i := 0; i < 10; i++ {
		connect()
	}
}
