package tasks

import "testing"

func TestExecConfig_Serialize(t *testing.T) {
	type fields struct {
		Command string
		Args    []string
		Timeout int
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		e := &ExecConfig{
			Command: tt.fields.Command,
			Args:    tt.fields.Args,
			Timeout: tt.fields.Timeout,
		}
		got, err := e.Serialize()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. ExecConfig.Serialize() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. ExecConfig.Serialize() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestExec(t *testing.T) {
	type args struct {
		config []interface{}
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
		got, err := Exec(tt.args.config...)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Exec() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. Exec() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
