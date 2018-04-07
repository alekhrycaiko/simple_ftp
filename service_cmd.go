package main

/*
* Handles the FTP service commands as defined by RFC 959
 */
import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func handleRetr(c *client) {
	if c.pasv == nil {
		sendMessage(c, 425)
		return
	}
	defer c.pasv.Close()
	if len(c.input) < 2 {
		sendMessage(c, 500)
		return
	}
	sendMessage(c, 150)
	// TODO: deal with encoding and type specifics.
	fp := filepath.Join(append(c.path, c.input[1])...)
	fd, err := os.Open(fp)
	if err != nil {
		sendMessage(c, 550)
		return
	}
	defer fd.Close()
	buff := make([]byte, 512)
	for {
		_, err := fd.Read(buff)
		if err == io.EOF {
			break
		}
		c.pasv.Write(buff)
	}
	sendMessage(c, 226)

}

/**
* Sends list of dir and files in present folder over pasv connection.
 */
func handleNlst(c *client) {
	if c.pasv == nil {
		sendMessage(c, 425)
		return
	}
	defer c.pasv.Close()
	fp := filepath.Join(c.path...)
	fmt.Println(fp)
	files, _ := ioutil.ReadDir(fp)
	writer := bufio.NewWriter(c.pasv)
	sendMessage(c, 150)
	for _, file := range files {
		s := []byte(fmt.Sprintf("%s %d\r\n", file.Name(), file.Size()))
		_, err := writer.Write(s)
		if err != nil {
			// TODO: check error code.
			sendMessage(c, 500)
			break
		}
	}
	writer.Flush()
	sendMessage(c, 226)
}
