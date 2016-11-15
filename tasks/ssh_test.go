package tasks

import (
	"reflect"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestSSH(t *testing.T) {
	type args struct {
		config []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := SSH(tt.args.config...)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. SSH() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. SSH() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSSHConfig_Serialize(t *testing.T) {
	tests := []struct {
		name    string
		s       *SSHConfig
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.s.Serialize()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. SSHConfig.Serialize() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. SSHConfig.Serialize() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSSHConfig_GetKeyAuths(t *testing.T) {
	tests := []struct {
		name      string
		s         *SSHConfig
		wantAuths []ssh.AuthMethod
		wantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotAuths, err := tt.s.GetKeyAuths()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. SSHConfig.GetKeyAuths() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotAuths, tt.wantAuths) {
			t.Errorf("%q. SSHConfig.GetKeyAuths() = %v, want %v", tt.name, gotAuths, tt.wantAuths)
		}
	}
}

func Test_parseKeyFiles(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name      string
		args      args
		wantAuths []ssh.AuthMethod
		wantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotAuths, err := parseKeyFiles(tt.args.paths)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. parseKeyFiles() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotAuths, tt.wantAuths) {
			t.Errorf("%q. parseKeyFiles() = %v, want %v", tt.name, gotAuths, tt.wantAuths)
		}
	}
}

func Test_loadEnvAgent(t *testing.T) {
	tests := []struct {
		name      string
		wantAuths []ssh.AuthMethod
		wantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotAuths, err := loadEnvAgent()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. loadEnvAgent() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotAuths, tt.wantAuths) {
			t.Errorf("%q. loadEnvAgent() = %v, want %v", tt.name, gotAuths, tt.wantAuths)
		}
	}
}

func Test_loadDefaultKeys(t *testing.T) {
	tests := []struct {
		name      string
		wantAuths []ssh.AuthMethod
		wantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotAuths, err := loadDefaultKeys()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. loadDefaultKeys() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotAuths, tt.wantAuths) {
			t.Errorf("%q. loadDefaultKeys() = %v, want %v", tt.name, gotAuths, tt.wantAuths)
		}
	}
}

func Test_fileExists(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := fileExists(tt.args.name); got != tt.want {
			t.Errorf("%q. fileExists() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
