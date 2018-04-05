package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Parses the incoming ftp connections command and returns
// the command line argument and its value.
func parseConnInput(input string) []string {
	fmt.Println("P")
	reader := bufio.NewReader(strings.NewReader(strings.Trim(input, " ")))
	str, _ := reader.ReadString('\n')
	return strings.Split(str, " ")
}

// Handles input from new connection.
func handleNewConnection(conn net.Conn) {
	client := &client{
		conn:   conn,
		writer: bufio.NewWriter(conn),
		path:   []string{"."},
		login:  false,
		ftype:  'A',
		mode:   'S',
		input:  []string{},
	}
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	defer conn.Close()
	sendMessage(client, 220)
	scanner := bufio.NewScanner(conn)
	client.scanner = scanner
	for scanner.Scan() {
		fmt.Println("T")
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:\r\n", err)
			break
		}
		text := scanner.Text()
		if len(text) == 0 {
			sendMessage(client, 500)
			continue
		}
		cmds := parseConnInput(text)
		client.input = cmds
		val, ok := cmdMap[cmds[0]]
		if ok {
			fmt.Println("O")
			val(client)
		} else {
			fmt.Println("F")
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
