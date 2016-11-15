package operator

import (
	"testing"

	"github.com/olebedev/config"
)

func TestConfigLogging(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		ConfigLogging(tt.args.cfg)
	}
}
