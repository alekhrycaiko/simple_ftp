package main

import (
	"fmt"
	"os"
	"path"
)

func handleAccessCommandUser(c *client) {
	if len(c.input) >= 2 && c.input[1] == "anonymous" {
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

/**
* Changes the working directory.
 */
func handleAccessCommandCwd(c *client) {
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
