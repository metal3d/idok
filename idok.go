package main

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var PORT string
var HOST string
var LISTEN string
var RASPIP string
var SSHPORT int

// message to send to stop media
const stopbody = `{"id":1,"jsonrpc":"2.0","method":"Player.Stop","params":{"playerid": %d}}`

// get player id
const getplayer = `{"id":1, "jsonrpc":"2.0","method":"Player.GetActivePlayers"}`

// the message to lauch local media
const body = `{
	"id":1,"jsonrpc":"2.0",
	"method":"Player.Open",
	"params": {
		"item": {
		   "file": "%s"
		 }
	 }
 }`

// response of get players
type itemresp struct {
	Id      int
	Jsonrpc string
	Result  []map[string]interface{}
}

// return active player from XBMC
func getActivePlayer() *itemresp {
	r, _ := http.Post(HOST, "application/json", bytes.NewBufferString(getplayer))
	response, _ := ioutil.ReadAll(r.Body)
	resp := &itemresp{}
	resp.Result = make([]map[string]interface{}, 0)
	json.Unmarshal(response, resp)
	return resp
}

// test if media is playing, if not then quit
func checkPlaying() {

	tick := time.Tick(3 * time.Second)
	for _ = range tick {
		resp := getActivePlayer()
		if len(resp.Result) == 0 {
			os.Exit(0)
		}
	}

}

// when quiting (CTRL+C for example) - tell to XBMC to stop
func onQuit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-c:
		fmt.Println("Quiting")
		resp := getActivePlayer()
		var playerid int
		for _, result := range resp.Result {
			for key, val := range result {
				if key == "playerid" {
					playerid = int(val.(float64))
				}
			}
		}

		http.Post(HOST, "application/json", bytes.NewBufferString(fmt.Sprintf(stopbody, playerid)))
		os.Exit(0)
	}
}

// begin to locally listen http to serve media
func send(host, file string, port int) {

	u := url.URL{Path: file}
	file = u.String()
	//_body := fmt.Sprintf(body, "http://"+LISTEN+":"+PORT+"/"+file)
	addr := fmt.Sprintf("http://%s:%d/%s", host, port, file)
	_body := fmt.Sprintf(body, addr)

	r, err := http.Post(HOST, "application/json", bytes.NewBufferString(_body))
	if err != nil {
		log.Fatal(err)
	}
	response, _ := ioutil.ReadAll(r.Body)
	log.Println(string(response))
}

// return local ip that matches kodi network
// ignoring loopback and other net interfaces
func getLocalInterfaceIP() (string, error) {
	ips, _ := net.LookupIP(RASPIP)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Error while checking you interfaces: %v", err)
	}
	for _, ip := range ips {
		mask := ip.DefaultMask()
		for _, iface := range ifaces {
			if iface.Flags&net.FlagLoopback != 0 {
				continue
			}

			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					if v.Mask.String() == mask.String() {
						return v.IP.String(), nil
					}
				}

			}
		}
	}
	return "", errors.New("Unable to get local ip")
}

// open a port locally and tell to kodi to stream
// from this port
func httpserve(file, dir string, port int) {

	localip, err := getLocalInterfaceIP()
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
	go send(localip, file, port)

	// handle CTRL+C to stop
	go onQuit()

	// and wait media end
	go checkPlaying()

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}

// Dig tunnel to kodi, open a port and bind socket to
// the local http server
func sshforward(config *ssh.ClientConfig, file, dir string) {

	// Setup sshClientConn (type *ssh.ClientConn)
	sshClientConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", RASPIP, SSHPORT), config)
	if err != nil {
		log.Fatal(err)
	}

	// Setup sshConn (type net.Conn)
	// Because dropbear doesn't accept :0 port to open random port
	// we do the randomisation ourself
	rand.Seed(int64(time.Now().Nanosecond()))
	port := 10000 + rand.Intn(9999)
	tries := 0
	sshConn, err := sshClientConn.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	for err != nil && tries < 500 {
		port = 10000 + rand.Intn(9999)
		sshConn, err = sshClientConn.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		tries++
	}
	log.Println("Listening port on raspberry: ", port)

	// send xbmc the file query
	go send("localhost", file, port)
	// handle CTRL+C to stop
	go onQuit()
	// and wait media end
	go checkPlaying()

	// now serve file
	fullpath := filepath.Join(dir, file)
	http.Serve(sshConn, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fullpath)
	}))
}

func main() {

	// flags
	xbmcaddr := flag.String("target", "", "xbmc/kodi ip (raspbmc address, ip or hostname)")
	username := flag.String("login", "", "jsonrpc login (configured in xbmc settings)")
	password := flag.String("password", "", "jsonrpc password (configured in xbmc settings)")
	viassh := flag.Bool("ssh", false, "Use SSH Tunnelling (need ssh user and password)")
	port := flag.Int("port", 8080, "local port (ignored if you use ssh option)")
	sshuser := flag.String("sshuser", "pi", "ssh login")
	sshpassword := flag.String("sshpass", "raspberry", "ssh password")
	sshport := flag.Int("sshport", 22, "target ssh port")

	flag.Parse()

	if *xbmcaddr == "" {
		fmt.Println("You must provide the xbmc server address")
		os.Exit(1)
	}

	HOST = *xbmcaddr
	RASPIP = *xbmcaddr
	SSHPORT = *sshport

	// XBMC can be configured to have username/password
	if *username != "" {
		HOST = *username + ":" + *password + "@" + HOST
	}
	HOST = "http://" + HOST + "/jsonrpc"

	if len(flag.Args()) < 1 {
		fmt.Println("You must provide a file to serve")
		os.Exit(2)
	}

	// find the good path
	toserve := flag.Arg(0)
	dir := "."
	toserve, _ = filepath.Abs(toserve)
	file := filepath.Base(toserve)
	dir = filepath.Dir(toserve)

	if *viassh {
		config := &ssh.ClientConfig{
			User: *sshuser,
			Auth: []ssh.AuthMethod{
				ssh.Password(*sshpassword),
			},
		}

		// serve !
		sshforward(config, file, dir)
	} else {
		httpserve(file, dir, *port)
	}
}
