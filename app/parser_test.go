package main

import (
	"slices"
	"strings"
	"testing"
)

func TestParseRequest(t *testing.T) {
	req := strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
	wantCmd := "hello"
	wantArgs := []string{"world"}
	gotCmd, gotArgs, err := parseRequest(req)
	if err != nil {
		t.Fatalf("got unexpected err %v", err)
	}

	if gotCmd != "hello" {
		t.Fatalf("got %v, wanted %v", gotCmd, wantCmd)
	}

	if slices.Equal(gotArgs, wantArgs) {
		t.Fatalf("got %v, wanted %v", gotArgs, wantArgs)
	}
}
