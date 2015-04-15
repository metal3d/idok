// This package gives method to send streams as an http server.
package asserver

import (
	"fmt"
	"github.com/metal3d/idok/utils"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

// Open a port locally and tell to kodi to stream
// from this port
func HttpServe(file, dir string, port int) {

	localip, err := utils.GetLocalInterfaceIP()
	log.Println(localip)
	if err != nil {
		log.Fatal(err)
	}

	// handle file http response
	fullpath := filepath.Join(dir, file)
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fullpath)
	}))

	// send xbmc the file query
	go utils.Send("http", localip, file, port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil));
}

// Serve STDIN stream from a local port
func TCPServeStdin(port int) {

	localip, err := utils.GetLocalInterfaceIP()
	log.Println(localip)
	if err != nil {
		log.Fatal(err)
	}

	// send xbmc the file query
	go utils.Send("tcp", localip, "", port)
	con, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatal(err)
	}
	c, err := con.Accept()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(c, os.Stdin)

}
