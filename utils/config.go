package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
)

type Config struct {
	// target ip or hostname
	Target string

	// target jsonrpc port
	Targetport int

	// Kodi Username
	User string

	// Kodi password
	Password string

	// host url to jsonrpc
	JsonRPC string

	// Local port if computer should serve
	Localport int

	// SSH port
	Sshport int

	// SSH user
	Sshuser string

	// SSHPassord
	Sshpassword string

	// Use SSH to stream
	Ssh bool
}

var GlobalConfig *Config

// Set the target host, port and ssh jsonrpc user/pass
func SetTarget(conf *Config) {

	host := conf.Target
	// XBMC can be configured to have username/password
	if conf.User != "" {
		host = conf.User + ":" + conf.Password + "@" + conf.Target
	}

	conf.JsonRPC = fmt.Sprintf("http://%s:%d/jsonrpc", host, conf.Targetport)

	// assign package conf
	GlobalConfig = conf
}

// Try to get config files
func CheckLocalConfigFiles() (string, bool) {
	u, _ := user.Current()
	home := u.HomeDir

	filelist := []string{
		"./idok.conf",
		home + "/.config/idok/idok.conf",
		home + "/.local/etc/idok.conf",
		"/etc/idok.conf",
	}

	for _, file := range filelist {
		if _, err := os.Stat(file); err == nil {
			return file, true
		}
	}

	return "", false
}

func LoadLocalConfig(filename string, config *Config) {
	content, _ := ioutil.ReadFile(filename)
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if strings.Trim(line, "") == "" {
			continue
		}
		// comments
		if line[0] == '#' {
			continue
		}
		// Get key = value...
		val := strings.SplitN(line, "=", 2)
		for i, _ := range val {
			val[i] = strings.Trim(val[i], " ")
		}

		// val[0] = key, val[1] = value
		switch strings.ToLower(val[0]) {
		case "target":
			config.Target = val[1]
		case "targetport":
			if val[1] == "" {
				continue
			}
			port, err := strconv.Atoi(val[1])
			if err != nil {
				log.Fatal("Target port in config file should be integer")
			}
			config.Targetport = int(port)
		case "login":
			config.User = val[1]
		case "password":
			config.Password = val[1]
		case "localport":
			port, err := strconv.Atoi(val[1])
			if err != nil {
				log.Fatal("Local port in config file should be integer")
			}
			config.Localport = int(port)
		case "sshuser":
			config.Sshuser = val[1]
		case "sshpass":
			config.Sshpassword = val[1]
		case "sshport":
			if val[1] == "" {
				continue
			}
			port, err := strconv.Atoi(val[1])
			if err != nil {
				log.Fatal("SSH port in config file should be integer")
			}
			config.Sshport = int(port)
		case "ssh":
			if val[1] == "true" {
				config.Ssh = true
			}
		default:
			log.Fatalf("Bad Key in configuration file: %s\n", val)
		}

	}
}

func PrintExampleConfig() {

	fmt.Println(`# blank value means default
#
# Idok checks if that file exists in that order:
# - ./idok.conf
# - $HOME/.config/idok/idok.conf
# - $HOME/.local/etc/idok.conf
# - /etc/idok.conf
#
# If option is given, it override the configuration file corresponding option
#
# IP or hostname of Kodi/XBMC
# (-target)
target = 

# port to connect jsonrpc
# (-targetport)
targetport = 

# Kodi/XBMC jsonrpc username and password
# (-login -password)
login = 
password = 

# SSH user/password (user is needed for ssh even if you sent keypair)
# (-sshuser -sshpass)
sshuser =
sshpass =

# if you changed ssh port from 22 to other
# (-sshport)
sshport = 

# force ssh usage (true or false)
# (-ssh -nossh)
ssh = 
`)
}
