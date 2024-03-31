package main

import "testing"

func TestEncodeBulkString(t *testing.T) {
	want := []byte("$3\r\nhey\r\n")
	got := encodeBulkString("hey")
	if string(got) != string(want) {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestEncodeSimpleString(t *testing.T) {
	want := []byte("+OK\r\n")
	got := encodeSimpleString("OK")
	if string(got) != string(want) {
		t.Fatalf("expected %s, got %s", want, got)
	}
}
