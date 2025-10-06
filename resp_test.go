package main

import (
	"strings"
	"testing"
)

func TestRead_BulkString(t *testing.T) {
	data := "$5\r\nhello\r\n"
	r := NewResp(strings.NewReader(data))

	v, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.typ != "bulk" {
		t.Fatalf("expected type 'bulk', got %q", v.typ)
	}
	if v.bulk != "hello" {
		t.Fatalf("expected bulk value 'hello', got %q", v.bulk)
	}
}

func TestRead_Array(t *testing.T) {
	data := "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n"
	r := NewResp(strings.NewReader(data))

	v, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.typ != "array" {
		t.Fatalf("expected type 'array', got %q", v.typ)
	}
	if len(v.array) != 2 {
		t.Fatalf("expected array length 2, got %d", len(v.array))
	}
	if v.array[0].bulk != "get" || v.array[1].bulk != "key" {
		t.Fatalf("unexpected array contents: %+v", v.array)
	}
}
