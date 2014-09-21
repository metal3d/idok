package utils

import (
	"flag"
	"fmt"
	"os"
)

// Usage() prints command line documentation
func Usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: ")
	fmt.Fprintf(os.Stderr, "%s [options] mediafile|youtubeurl|streamurl\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Opening external URL dosen't open local or remote port. Your media center will fetch data itself.\n\n")
	fmt.Fprintf(os.Stderr, "You may be able to stream stdout -> stdin:")
	fmt.Fprintf(os.Stderr, "\n\t%s [options] -stdin < file\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Or:\n\tcommand | %s [options] -stdin \n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Using ssh option is only managed for local files.\n")
	fmt.Fprintf(os.Stderr, "Default mode is HTTP mode, it opens :8080 port on your host and send message to Kodi to read from that port. So, you must configure your firewall to open that port. You can override used port with -port option.\n")
	fmt.Fprintf(os.Stderr, "You can use SSH with -ssh option, %s will try to use key pair authtification, then use -sshpass to try login/password auth. With -ssh, you should change -sshuser if your Kodi user is not \"pi\" (default on raspbmc)\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "To be able to authenticate without password, use the command:\n\n\tssh-copy-id USER@KODI_HOST\n\nwhere USER is the Kodi user (pi) and KODI_HOST the ip or hostname of Kodi host.")
	fmt.Fprintf(os.Stderr, "\n\nOptions:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
You can create a configuration file in current directory, or:
- $HOME/.config/idok/idok.conf
- $HOME/.local/etc/idok.conf
- /etc/idok.conf

Configuration file will set default options. Using options in the command line will override configuration file

The -conf-example option prints a default configuration file.
`)
	fmt.Fprintf(os.Stderr, "\n")
}
