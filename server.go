package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Parses the incoming ftp connections command and returns
// the command line argument and its value.
func parseConnInput(input string, num int) (string, error) {
	fmt.Printf("INPUT  IS %s", input)
	if len(input) == 0 {
		return "", errors.New("No arg")
	}
	words := strings.Fields(input)
	s := strings.ToLower(words[num])
	return s, nil
}

// Handles input from new connection.
func handleNewConnection(conn net.Conn) {
	client := &client{
		conn:   conn,
		writer: bufio.NewWriter(conn),
		dir:    "./file_directory",
		login:  false,
	}
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	defer conn.Close()
	sendMessage(client, 220)
	scanner := bufio.NewScanner(conn)
	client.scanner = scanner
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:\r\n", err)
			break
		}
		text := scanner.Text()
		cmd, _ := parseConnInput(text, 0)
		fmt.Printf("%s\r\n", cmd)
		if len(text) == 0 {
			sendMessage(client, 500)
			continue
		}
		val, ok := cmdMap[cmd]
		if ok {
			val(client)
		} else {
			sendMessage(client, 502)
		}
	}
	fmt.Printf("Scanner decided to close.")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Invalid argument count")
		os.Exit(1)
	}
	port := ":" + os.Args[1]
	address, err := net.ResolveTCPAddr("tcp4", port)
	if err != nil {
		fmt.Printf("Error %s", err)
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", address)
	defer listener.Close()
	fmt.Printf("FTP server listening at port %s", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		if conn != nil {
			go handleNewConnection(conn)
		}
	}
}
