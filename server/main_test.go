package main

import "testing"

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for range tests {
		main()
	}
}

func Test_getWorkerName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := getWorkerName(); got != tt.want {
			t.Errorf("%q. getWorkerName() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
