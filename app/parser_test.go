package main

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestParseHelloWorldArray(t *testing.T) {
	req := strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
	want := []any{"hello", "world"}
	got, err := parseRequest(req)
	if err != nil {
		t.Fatalf("got unexpected err %v", err)
	}

	v := reflect.ValueOf(got)
	if v.Kind() != reflect.Slice {
		t.Fatalf("got unexpected type %T", t)
	}

	if !slices.Equal(got, want) {
		t.Fatalf(`parsing %q should return %#v, got %#v`, req, want, got)
	}
}
