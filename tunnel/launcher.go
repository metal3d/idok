package tunnel

import (
	"fmt"
	"github.com/metal3d/idok/tunnel/go.crypto/ssh"
	"github.com/metal3d/idok/utils"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// SshForward digs a tunnel to xbmc/kodi, then open a port and bind socket to
// the local http server
func SshForward(config *ssh.ClientConfig, file, dir string) {

	// Setup sshClientConn (type *ssh.ClientConn)
	sshClientConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", utils.GlobalConfig.Target, utils.GlobalConfig.Sshport), config)
	if err != nil {
		log.Fatal(err)
	}

	// Setup sshConn (type net.Conn)
	// Because dropbear doesn't accept :0 port to open random port
	// we do the randomisation ourself
	rand.Seed(int64(time.Now().Nanosecond()))
	dport := 10000 + rand.Intn(9999)
	tries := 0
	sshConn, err := sshClientConn.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", dport))
	for err != nil && tries < 500 {
		dport = 10000 + rand.Intn(9999)
		sshConn, err = sshClientConn.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", dport))
		tries++
	}
	log.Println("Listening port on raspberry: ", dport)

	// send xbmc the file query
	go utils.Send("http", "localhost", file, dport)
	// handle CTRL+C to stop
	go utils.OnQuit()

	// now serve file
	fullpath := filepath.Join(dir, file)
	http.Serve(sshConn, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fullpath)
	}))
}

// SshForwardStdin reads stdin and stream this to distant socket
// through SSH tunnel
func SshForwardStdin(config *ssh.ClientConfig) {

	// Setup sshClientConn (type *ssh.ClientConn)
	sshClientConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", utils.GlobalConfig.Target, utils.GlobalConfig.Sshport), config)
	if err != nil {
		log.Fatal(err)
	}

	// Setup sshConn (type net.Conn)
	// Because dropbear doesn't accept :0 port to open random port
	// we do the randomisation ourself
	rand.Seed(int64(time.Now().Nanosecond()))
	dport := 10000 + rand.Intn(9999)
	tries := 0
	sshConn, err := sshClientConn.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", dport))
	for err != nil && tries < 500 {
		dport = 10000 + rand.Intn(9999)
		sshConn, err = sshClientConn.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", dport))
		tries++
	}
	log.Println("Listening port on the target: ", dport)

	// send xbmc the file query
	go utils.Send("tcp", "127.0.0.1", "", dport)
	if err != nil {
		log.Fatal(err)
	}
	c, err := sshConn.Accept()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(c, os.Stdin)

	// handle CTRL+C to stop
	utils.OnQuit()
}
