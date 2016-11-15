package operator

import (
	"reflect"
	"testing"

	"github.com/olebedev/config"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantCfg *config.Config
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotCfg, err := GetConfig()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. GetConfig() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
			t.Errorf("%q. GetConfig() = %v, want %v", tt.name, gotCfg, tt.wantCfg)
		}
	}
}

func TestPrettyPrintFlagMap(t *testing.T) {
	type args struct {
		m      map[string]interface{}
		prefix []string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		PrettyPrintFlagMap(tt.args.m, tt.args.prefix)
	}
}

func Test_combineConfigs(t *testing.T) {
	type args struct {
		cfgs []*config.Config
	}
	tests := []struct {
		name  string
		args  args
		wantR *config.Config
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if gotR := combineConfigs(tt.args.cfgs...); !reflect.DeepEqual(gotR, tt.wantR) {
			t.Errorf("%q. combineConfigs() = %v, want %v", tt.name, gotR, tt.wantR)
		}
	}
}

func Test_isCfgFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := isCfgFile(tt.args.path); got != tt.want {
			t.Errorf("%q. isCfgFile() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
