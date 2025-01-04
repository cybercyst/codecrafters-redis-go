package server

import (
	"bufio"
	"slices"
	"strings"
	"testing"
)

func TestParseRequest(t *testing.T) {
	req := strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
	scanner := bufio.NewScanner(req)
	wantArgs := []string{"hello", "world"}
	gotArgs, err := parseRESP(scanner)
	if err != nil {
		t.Fatalf("got unexpected err %v", err)
	}

	if !slices.Equal(gotArgs, wantArgs) {
		t.Fatalf("got %v, wanted %v", gotArgs, wantArgs)
	}
}
