package utils

import (
	"fmt"
)

// the target host address (http://...)
var HOST string

// The target IP
var TARGETIP string

// ssh port number (default is 22)
var SSHPORT int

// Set the target host, port and ssh jsonrpc user/pass
func SetTarget(host string, port int, username, password *string) {
	HOST = host
	TARGETIP = host

	// XBMC can be configured to have username/password
	if *username != "" {
		HOST = *username + ":" + *password + "@" + HOST
	}

	HOST = fmt.Sprintf("http://%s:%d/jsonrpc", HOST, port)
}
