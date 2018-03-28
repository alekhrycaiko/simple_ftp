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

func handleError(err error) {
	fmt.Printf("Error %s", err)
	os.Exit(1)
}

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

type client struct {
	c    net.Conn
	dir  string
	path string
}

// Main handling function for new connections. Handles responses and reads.
func handleNewConnection(conn net.Conn, cmdMap map[string]interface{}) {
	conn.Write(formatMsg(220, "Accepted Connection to FTP. Enjoy!"))
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	loggedIn := false
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:\r\n", err)
			break
		}
		text := scanner.Text()
		cmd, _ := parseConnInput(text, 0)
		fmt.Printf("%s\r\n", cmd)

		if cmd == "quit" {
			handleAccessCommandQuit(conn)
			return
		}

		if len(text) == 0 {
			conn.Write([]byte("Invalid command. \r\n"))
			continue
		}

		if !loggedIn {
			switch cmd {
			case "user":
				if len(text) > 5 {
					arg, err := parseConnInput(text, 1)
					if err != nil {
						conn.Write(formatMsg(331, "FTP Server is Anonymous only"))
						break
					}
					handleAccessCommandUser(conn, arg, loggedIn)
				} else {
					conn.Write(formatMsg(331, "FTP Server is Anonymous only"))
				}
			case "pass":
				handleAccessCommandPass(conn, &loggedIn)
			default:
				conn.Write(formatMsg(530, "Please login with USER and Pass"))
			}
		} else {
			switch cmd {
			case "pasv":
				setupPasvConnection(conn)
			default:
				conn.Write(formatMsg(500, "Unsupported command"))
			}
		}
	}
	fmt.Printf("Scanner decided to close.")
}

func main() {
	if len(os.Args) != 2 {
		handleError(errors.New("Invalid argument count"))
	}
	port := ":" + os.Args[1]
	address, err := net.ResolveTCPAddr("tcp4", port)
	if err != nil {
		handleError(err)
	}
	listener, err := net.ListenTCP("tcp", address)

	defer listener.Close()
	fmt.Printf("FTP server listening at port %s", port)
	cmdMap := map[string]interface{}{
		"quit": handleAccessCommandQuit{},
		"user": handleAccessCommandUser{},
		"pass": handleAccessCommandPass{},
		"pasv": setupPasvConnection{},
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		if conn != nil {
			go handleNewConnection(conn, cmdMap)
		}
	}
}
