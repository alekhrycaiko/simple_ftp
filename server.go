package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// parseConnInput parses the incoming ftp connections command and returns
// the command line arguments.
func parseConnInput(input string) []string {
	reader := bufio.NewReader(strings.NewReader(strings.Trim(input, " ")))
	str, _ := reader.ReadString('\n')
	return strings.Split(str, " ")
}

// handleNewConnection is the primary handler (main) for a new connection.
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
		if ok && client.login || (ok && !client.login && (cmds[0] == "user" || cmds[0] == "pass")) {
			val(client)
		} else if ok && !client.login {
			sendMessage(client, 331)
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
	num := flag.String("port", "3000", "Port used")
	flag.Parse()
	port := ":" + *num
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
