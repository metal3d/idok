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
	"os/user"
	"path/filepath"
	"syscall"
	"time"
)

var PORT string
var HOST string
var LISTEN string
var RASPIP string
var SSHPORT int

const (
	VERSION = "0.2.2"

	// message to send to stop media
	stopbody = `{"id":1,"jsonrpc":"2.0","method":"Player.Stop","params":{"playerid": %d}}`

	// get player id
	getplayer = `{"id":1, "jsonrpc":"2.0","method":"Player.GetActivePlayers"}`

	// the message to lauch local media
	body = `{
	"id":1,"jsonrpc":"2.0",
	"method":"Player.Open",
	"params": {
		"item": {
		   "file": "%s"
		 }
	 }
 }`
)

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

// Parse local ssh private key to get signer
func parseSSHKeys() ssh.Signer {
	u, _ := user.Current()
	home := u.HomeDir
	id_rsa_priv := filepath.Join(home, ".ssh", "id_rsa")
	content, err := ioutil.ReadFile(id_rsa_priv)
	if err != nil {
		log.Println("no id_rsa key found")
		return nil
	}

	private, err := ssh.ParsePrivateKey(content)
	if err != nil {
		log.Println("Unable to parse private key")
		return nil
	}
	fmt.Println(private)
	return private

}

func main() {

	// flags
	xbmcaddr := flag.String("target", "", "xbmc/kodi ip (raspbmc address, ip or hostname)")
	username := flag.String("login", "", "jsonrpc login (configured in xbmc settings)")
	password := flag.String("password", "", "jsonrpc password (configured in xbmc settings)")
	viassh := flag.Bool("ssh", false, "Use SSH Tunnelling (need ssh user and password)")
	port := flag.Int("port", 8080, "local port (ignored if you use ssh option)")
	sshuser := flag.String("sshuser", "pi", "ssh login")
	sshpassword := flag.String("sshpass", "", "ssh password")
	sshport := flag.Int("sshport", 22, "target ssh port")
	version := flag.Bool("version", false, fmt.Sprintf("Print the current version (%s)", VERSION))

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n%s [options] mediafile\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Default mode is HTTP mode, it opens :8080 port on your host and send message to Kodi to open that port.\n")
		fmt.Fprintf(os.Stderr, "You can use SSH with -ssh option, %s will try to use key pair authtification, then use -sshpass to try login/password auth. With -ssh, you should change -sshuser if your Kodi user is not \"pi\" (default on raspbmc)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "To be able to authenticate without password, use the command:\n\n\tssh-copy-id USER@KODI_HOST\n\nwhere USER is the Kodi user (pi) and KODI_HOST the ip or hostname of Kodi host.")
		fmt.Fprintf(os.Stderr, "\n\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if *xbmcaddr == "" {
		fmt.Println("\033[33mYou must provide the xbmc server address\033[0m")
		flag.Usage()
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
		fmt.Println("\033[33mYou must provide a file to serve\033[0m")
		flag.Usage()
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
				ssh.PublicKeys(parseSSHKeys()),
				ssh.Password(*sshpassword),
			},
		}

		// serve !
		sshforward(config, file, dir)
	} else {
		httpserve(file, dir, *port)
	}
}
