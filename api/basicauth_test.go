package api

import (
	"net/http"
	"reflect"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	type args struct {
		pw      string
		handler http.HandlerFunc
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := BasicAuth(tt.args.pw, tt.args.handler); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. BasicAuth() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestValidate(t *testing.T) {
	type args struct {
		pw   string
		test string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := Validate(tt.args.pw, tt.args.test); got != tt.want {
			t.Errorf("%q. Validate() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
