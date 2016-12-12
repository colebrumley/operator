package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"time"

	"bytes"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSH dispatches a command to a remote node via SSH
func SSH(config ...interface{}) (string, error) {
	var (
		cfg         SSHConfig
		auths       []ssh.AuthMethod
		err         error
		currentUser *user.User
		rawData     []byte
	)

	if len(config) > 0 {
		rawData = []byte(config[0].(string))
	} else {
		return "", fmt.Errorf("No SSH config provided!")
	}

	// Convert the config back into a struct from JSON
	if err := json.Unmarshal(rawData, &cfg); err != nil {
		return "", err
	}

	// Make sure we have the necessary info to actually do something
	if !cfg.Validate() {
		return "", fmt.Errorf("Insufficient information provided to establish an SSH connection")
	}

	// Set the User to the currect user if it's undefined
	if len(cfg.User) == 0 {
		if currentUser, err = user.Current(); err != nil {
			return "", err
		}
		cfg.User = currentUser.Name
	}

	// Get the auth methods from provided creds. Prefer keys, fall back to password if no keys are defined.
	if len(cfg.Keyfiles) > 0 {
		log.Debugln("Loading SSH keyfiles", cfg.Keyfiles)
		auths, err = GetKeyAuths(cfg.Keyfiles...)
	} else if len(cfg.Password) > 0 {
		log.Debugln("Using password for SSH access")
		auths = append(auths, ssh.Password(cfg.Password))
	} else {
		log.Warningln("No SSH auth mechanisms defined, using defaults")
		auths, err = GetDefaultAuth()
	}
	if err != nil {
		return "", err
	}

	// Run the command and return the stringified output. Default to 10s connection timeout
	client, err := ssh.Dial("tcp", portAddrCheck(cfg.Host), &ssh.ClientConfig{
		User:    cfg.User,
		Auth:    auths,
		Timeout: 10 * time.Second,
	})
	defer client.Close()
	if err != nil {
		return "", err
	}
	session, err := client.NewSession()
	defer session.Close()
	if err != nil {
		return "", err
	}

	out := bytes.NewBuffer([]byte{})
	session.Stderr = out
	session.Stdout = out
	if cfg.Shell {
		session.Stdin = bytes.NewBufferString(cfg.Command)
		if err = session.Shell(); err != nil {
			return out.String(), err
		}
	} else {
		if err = session.Start(cfg.Command); err != nil {
			return out.String(), err
		}
	}

	err = session.Wait()
	return out.String(), err
}

// SSHConfig is the wrapper object for SSH settings to be passed through to the SSH task
type SSHConfig struct {
	User, Password, Host, Command string
	Shell                         bool
	Keyfiles                      []string
}

// Validate makes sure we have enough information to make a call
func (s *SSHConfig) Validate() bool {
	// Make sure we have the necessary info to actually do something
	if len(s.Command) == 0 || len(s.Host) == 0 || len(s.User) == 0 {
		return false
	}
	return true
}

// GetDefaultAuth is the default method for loading keys. Uses SSH_AUTH_SOCK
// if available, otherwise attempts to load default keys at ~/.ssh/id_rsa
// for Linux and ~/ssh/id_rsa (for windows, but it works on linux too).
func GetDefaultAuth() (auths []ssh.AuthMethod, err error) {
	if len(os.Getenv("SSH_AUTH_SOCK")) > 0 {
		auths, err = loadEnvAgent()
	} else {
		var currentUser *user.User
		if currentUser, err = user.Current(); err != nil {
			return
		}
		auths, err = GetKeyAuths(
			filepath.FromSlash(currentUser.HomeDir+"/.ssh/id_rsa"),
			filepath.FromSlash(currentUser.HomeDir+"/ssh/id_rsa"),
		)
	}
	return
}

// GetKeyAuths returns ssh.AuthMethod objects loaded from the provided SSH key paths.
func GetKeyAuths(keyfiles ...string) (auths []ssh.AuthMethod, err error) {
	for _, path := range keyfiles {
		var (
			pemBytes []byte
			signer   ssh.Signer
		)
		if !fileExists(path) {
			err = errors.New("Specified key does not exist")
			return
		}
		if pemBytes, err = ioutil.ReadFile(path); err != nil {
			return
		}
		if signer, err = ssh.ParsePrivateKey(pemBytes); err != nil {
			return
		}
		auths = append(auths, ssh.PublicKeys(signer))
	}
	return
}

func loadEnvAgent() (auths []ssh.AuthMethod, err error) {
	sshAuthSock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return
	}
	defer sshAuthSock.Close()
	ag := agent.NewClient(sshAuthSock)
	auths = []ssh.AuthMethod{ssh.PublicKeysCallback(ag.Signers)}
	return
}

// fileExists returns a bool if os.Stat returns an IsNotExist error
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func portAddrCheck(addr string) string {
	if len(strings.Split(addr, ":")) == 1 {
		return addr + ":22"
	}
	return addr
}
