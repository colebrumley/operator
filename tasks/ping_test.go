package tasks

import "testing"

func TestPing(t *testing.T) {
	type args struct {
		in0 []interface{}
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
		got, err := Ping(tt.args.in0...)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Ping() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. Ping() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
