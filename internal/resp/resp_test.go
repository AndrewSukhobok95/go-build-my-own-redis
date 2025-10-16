package resp

import (
	"strings"
	"testing"
)

func TestMarshal_BulkString(t *testing.T) {
	v := Value{typ: "bulk", bulk: "hello"}
	want := "$5\r\nhello\r\n"
	if string(v.marshalBulk()) != want {
		t.Errorf("got %q, want %q", string(v.marshalBulk()), want)
	}
}

func TestMarshal_Array(t *testing.T) {
	v := Value{typ: "array", array: []Value{
		{typ: "bulk", bulk: "SET"},
		{typ: "bulk", bulk: "key"},
		{typ: "bulk", bulk: "value"},
	}}
	want := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	if string(v.marshalArray()) != want {
		t.Errorf("got %q, want %q", string(v.marshalArray()), want)
	}
}

func TestMarshal(t *testing.T) {
	v := Value{
		typ: "array",
		array: []Value{
			{typ: "bulk", bulk: "foo"},
			{
				typ: "array",
				array: []Value{
					{typ: "bulk", bulk: "bar"},
				},
			},
		},
	}
	got := v.Marshal()
	want := "*2\r\n$3\r\nfoo\r\n*1\r\n$3\r\nbar\r\n"
	if string(got) != want {
		t.Errorf("Marshal() = %q, want %q", string(got), want)
	}
}

func TestRead_BulkString(t *testing.T) {
	data := "$5\r\nhello\r\n"
	r := NewReader(strings.NewReader(data))

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
	r := NewReader(strings.NewReader(data))

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
