package api

import (
	"net/http"
	"reflect"
	"testing"

	machinery "github.com/RichardKnop/machinery/v1"
)

func TestOperatorAPI_Start(t *testing.T) {
	type args struct {
		server *machinery.Server
	}
	tests := []struct {
		name string
		o    *OperatorAPI
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.o.Start(tt.args.server)
	}
}

func TestDefaultHandler(t *testing.T) {
	type args struct {
		rw  http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		DefaultHandler(tt.args.rw, tt.args.req)
	}
}

func Test_makeHandler(t *testing.T) {
	type args struct {
		name   string
		server *machinery.Server
	}
	tests := []struct {
		name string
		args args
		want func(http.ResponseWriter, *http.Request)
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := makeHandler(tt.args.name, tt.args.server); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. makeHandler() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
