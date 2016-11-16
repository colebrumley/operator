package tasks

import (
	"encoding/json"
	"errors"
	"os/exec"
	"time"

	"bytes"

	"fmt"

	"math"

	log "github.com/Sirupsen/logrus"
)

// ExecConfig is the wrapper object for Exec settings to be passed through to the Exec task
type ExecConfig struct {
	Command string
	Args    []string
	Timeout int
}

// Serialize returns a stringified ExecConfig
func (e *ExecConfig) Serialize() (string, error) {
	bytes, err := json.Marshal(e)
	return string(bytes), err
}

// Exec runs an external application or script. TODO: add a timeout
func Exec(config ...string) (string, error) {
	var cfg ExecConfig

	data := []byte(config[0])

	// Convert the config back into a struct from JSON
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}

	if len(cfg.Command) < 1 {
		err := errors.New("No command specified")
		return "", err
	}

	cmd := exec.Command(cfg.Command, cfg.Args...)
	buf := bytes.NewBuffer([]byte{})
	cmd.Stdout = buf
	cmd.Stderr = buf
	if err := cmd.Start(); err != nil {
		return "", err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	if cfg.Timeout == 0 {
		cfg.Timeout = math.MaxInt32
	}

	select {
	case <-time.After(time.Duration(cfg.Timeout) * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			log.Error("Failed to kill timed-out exec task: ", err)
			return "", err
		}
		log.Warn("Exec process timed out")
		return buf.String(), fmt.Errorf("Operation timed out after %v seconds", cfg.Timeout)
	case err := <-done:
		if err != nil {
			log.Warnf("Exec process completed with error = %v", err)
			return buf.String(), err
		}
	}
	log.Info("Exec process completed")
	return buf.String(), nil
}
