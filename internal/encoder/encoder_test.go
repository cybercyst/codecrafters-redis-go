package encoder

import "testing"

func TestEncodeEmptyArray(t *testing.T) {
	want := []byte("*0\r\n")
	got := EncodeArrayBulkString([]string{})
	if string(got) != string(want) {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestEncodeArray(t *testing.T) {
	want := []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
	got := EncodeArrayBulkString([]string{"hello", "world"})
	if string(got) != string(want) {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestEncodeBulkString(t *testing.T) {
	want := []byte("$3\r\nhey\r\n")
	got := EncodeBulkString("hey")
	if string(got) != string(want) {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestEncodeSimpleString(t *testing.T) {
	want := []byte("+OK\r\n")
	got := EncodeSimpleString("OK")
	if string(got) != string(want) {
		t.Fatalf("expected %s, got %s", want, got)
	}
}
