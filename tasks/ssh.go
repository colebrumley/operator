package tasks

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSH dispatches a command to a remote node via SSH
func SSH(config ...string) (string, error) {
	var cfg SSHConfig

	data := []byte(config[0])

	// Convert the config back into a struct from JSON
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}

	// Make sure we have the necessary info to actually do something
	if len(cfg.Command) == 0 {
		return "No command provided, nothing to do.", nil
	}

	if len(cfg.Host) == 0 {
		return "", errors.New("No host defined")
	}

	if len(cfg.User) == 0 {
		return "", errors.New("No user defined")
	}

	if len(cfg.Keyfiles) == 0 && len(cfg.Password) == 0 {
		return "", errors.New("No RSA keys or password defined, cannot authenticate")
	}

	// Get the auth methods from provided creds. Prefer keys, fall back to password if no keys are defined.
	auths, err := cfg.GetKeyAuths()
	if err != nil {
		return "", err
	}
	if len(auths) == 0 {
		if len(cfg.Password) == 0 {
			return "", errors.New("No auth methods defined, cannot continue!")
		}
		auths = append(auths, ssh.Password(cfg.Password))
	}

	// Run the command and return the stringified output
	client, err := ssh.Dial("tcp", portAddrCheck(cfg.Host), &ssh.ClientConfig{
		User: cfg.User,
		Auth: auths,
	})
	if err != nil {
		return "", err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	out, err := session.CombinedOutput(cfg.Command)
	return string(out), err
}

// SSHConfig is the wrapper object for SSH settings to be passed through to the SSH task
type SSHConfig struct {
	User, Password, Host, Command string
	Keyfiles                      []string
}

// Serialize returns a stringified SSHConfig
func (s *SSHConfig) Serialize() (string, error) {
	bytes, err := json.Marshal(s)
	return string(bytes), err
}

// GetKeyAuths returns the ssh.AuthMethod objects loaded from the provided SSH keys
func (s *SSHConfig) GetKeyAuths() (auths []ssh.AuthMethod, err error) {
	auths = []ssh.AuthMethod{}
	if len(s.Keyfiles) == 0 {
		// If no keys are provided see if there's an ssh-agent running
		if len(os.Getenv("SSH_AUTH_SOCK")) > 0 {
			if auths, err = loadEnvAgent(); err != nil {
				return
			}
		} else {
			// Otherwise use ~/.ssh/id_rsa or ~/ssh/id_rsa (for windows, but
			// it works on linux too)
			if auths, err = loadDefaultKeys(); err != nil {
				return
			}
		}
	} else {
		// Append each provided key to auths
		auths, err = parseKeyFiles(s.Keyfiles)
	}
	if len(auths) == 0 {
		err = errors.New("No auths parsed from provided keys")
	}
	return
}

func parseKeyFiles(paths []string) (auths []ssh.AuthMethod, err error) {
	for _, key := range paths {
		var (
			pemBytes []byte
			signer   ssh.Signer
		)
		if !fileExists(key) {
			err = errors.New("Specified key does not exist")
			return
		}
		pemBytes, err = ioutil.ReadFile(key)
		if err != nil {
			return
		}
		signer, err = ssh.ParsePrivateKey(pemBytes)
		if err != nil {
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

func loadDefaultKeys() (auths []ssh.AuthMethod, err error) {
	k := ""
	currentUser, err := user.Current()
	defaultKeyPathA := filepath.FromSlash(currentUser.HomeDir + "/.ssh/id_rsa")
	defaultKeyPathB := filepath.FromSlash(currentUser.HomeDir + "/ssh/id_rsa")
	if fileExists(defaultKeyPathA) {
		k = defaultKeyPathA
	} else if fileExists(defaultKeyPathB) {
		k = defaultKeyPathB
	}
	if len(k) == 0 {
		err = errors.New("No key specified")
		return
	}
	pemBytes, err := ioutil.ReadFile(k)
	if err != nil {
		return
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return
	}
	auths = []ssh.AuthMethod{ssh.PublicKeys(signer)}
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
