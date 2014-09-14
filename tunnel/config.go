package tunnel

import (
	"code.google.com/p/go.crypto/ssh"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// NewConfig returns a ssh.Config pointer with 3 auth method if possible, rsa key pair,
// dsa keypair and user/pass
func NewConfig(username, password string) *ssh.ClientConfig {
	u, _ := user.Current()
	home := u.HomeDir
	id_rsa_priv := filepath.Join(home, ".ssh", "id_rsa")
	id_dsa_priv := filepath.Join(home, ".ssh", "id_dsa")

	auth := []ssh.AuthMethod{}

	// Try to parse keypair
	if _, err := os.Stat(id_rsa_priv); err == nil {
		if keypair, err := parseSSHKeys(id_rsa_priv); err == nil {
			log.Println("Added RSA key")
			auth = append(auth, ssh.PublicKeys(keypair))
		}
	}
	if _, err := os.Stat(id_dsa_priv); err == nil {
		if keypair, err := parseSSHKeys(id_dsa_priv); err == nil {
			log.Println("Added DSA key")
			auth = append(auth, ssh.PublicKeys(keypair))
		}
	}

	// add password method
	auth = append(auth, ssh.Password(password))

	// and set config
	return &ssh.ClientConfig{
		User: username,
		Auth: auth,
	}

}
