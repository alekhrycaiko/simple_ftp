package main

import (
	"fmt"
	"os"
	"path"
)

// handleUser will send a message back to the client based on input credentials.
func handleUser(c *client) {
	if len(c.input) >= 2 && c.input[1] == "anonymous" {
		c.login = true
		sendMessage(c, 331)
	} else {
		sendMessage(c, 530)
	}
}

// handlePass will always succeed for client currently since we support anonymous mode.
// no passwords are required for this server.
func handlePass(c *client) {
	sendMessage(c, 230)
}

// handleQuit will close the current tcp connection.
func handleQuit(c *client) {
	sendMessage(c, 221)
	c.conn.Close()
}

// handleCwd changes the file directory of the client by pushing the given argument
// if it produces a valid path.
func handleCwd(c *client) {
	if len(c.input) < 2 {
		sendMessage(c, 500)
		return
	}
	arg := c.input[1]
	if arg == "./" || arg == "../" {
		sendMessage(c, 550)
		return
	}
	if arg == ".." {
		handleCdup(c)
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
}

// handleCdup moves file directory up, if we're not already at the parent '.'.
func handleCdup(c *client) {
	if len(c.path) == 1 {
		sendMessage(c, 550)
	} else {
		c.path = c.path[:len(c.path)-1]
		fmt.Println(c.path)
		sendMessage(c, 250)
	}
}
