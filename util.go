package main

import (
	"fmt"
)

/**
* Formats and a message with a status code.
 */
func formatMsg(status int, message string) []byte {
	return []byte(fmt.Sprintf("%d %s \r\n", status, message))
}
