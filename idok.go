package main

import (
	"./asserver"
	"./tunnel"
	"./utils"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	VERSION = "v1-alpha1"
)

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
	xbmcport := flag.Int("targetport", 80, "XBMC/Kodi jsonrpc port")
	stdin := flag.Bool("stdin", false, "Read file from stdin")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n%s [options] mediafile|youtubeurl\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Opening youtube urls dosen't open local or remote port.\n")
		fmt.Fprintf(os.Stderr, "Using ssh option is only managed for local files.\n")
		fmt.Fprintf(os.Stderr, "Default mode is HTTP mode, it opens :8080 port on your host and send message to Kodi to open that port.\n")
		fmt.Fprintf(os.Stderr, "You can use SSH with -ssh option, %s will try to use key pair authtification, then use -sshpass to try login/password auth. With -ssh, you should change -sshuser if your Kodi user is not \"pi\" (default on raspbmc)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "To be able to authenticate without password, use the command:\n\n\tssh-copy-id USER@KODI_HOST\n\nwhere USER is the Kodi user (pi) and KODI_HOST the ip or hostname of Kodi host.")
		fmt.Fprintf(os.Stderr, "\n\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// print the current version
	if *version {
		fmt.Println(VERSION)
		fmt.Println("Compiled for", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if *xbmcaddr == "" {
		fmt.Println("\033[33mYou must provide the xbmc server address\033[0m")
		flag.Usage()
		os.Exit(1)
	}

	utils.SSHPORT = *sshport
	utils.SetTarget(*xbmcaddr, *xbmcport, username, password)

	var dir, file string

	// we don't use stdin, so we should check if scheme is file, youtube or other...
	if !*stdin {
		if len(flag.Args()) < 1 {
			fmt.Println("\033[33mYou must provide a file to serve\033[0m")
			flag.Usage()
			os.Exit(2)
		}

		if youtube, vid := utils.IsYoutubeURL(flag.Arg(0)); youtube {
			log.Println("Youtube video, using youtube addon from XBMC/Kodi")
			utils.PlayYoutube(vid)
			os.Exit(0)
		}

		if ok, local := utils.IsOtherScheme(flag.Arg(0)); ok {
			log.Println("\033[33mWarning, other scheme could be not supported by you Kodi/XBMC installation. If doesn't work, check addons and stream\033[0m")
			utils.SendBasicStream(flag.Arg(0), local)
			os.Exit(0)
		}

		// find the good path
		toserve := flag.Arg(0)
		dir = "."
		toserve, _ = filepath.Abs(toserve)
		file = filepath.Base(toserve)
		dir = filepath.Dir(toserve)

	}

	if *viassh {
		config := tunnel.NewConfig(*sshuser, *sshpassword)
		// serve ssh tunnel !
		if !*stdin {
			tunnel.SshForward(config, file, dir)
		} else {
			tunnel.SshForwardStdin(config)
		}
	} else {
		// serve local port !
		if !*stdin {
			asserver.HttpServe(file, dir, *port)
		} else {
			asserver.HttpServeStdin(*port)
		}
	}
}
