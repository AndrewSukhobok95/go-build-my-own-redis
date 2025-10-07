package main

import "testing"

func TestHandlePing(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want Value
	}{
		{
			name: "PING with no args",
			args: nil,
			want: Value{typ: "string", str: "PONG"},
		},
		{
			name: "PING with one arg",
			args: []string{"hello"},
			want: Value{typ: "bulk", bulk: "hello"},
		},
		{
			name: "PING with too many args",
			args: []string{"foo", "bar"},
			want: Value{typ: "error", str: "ERR wrong number of arguments"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handlePing(tt.args)
			if got.typ != tt.want.typ || got.str != tt.want.str || got.bulk != tt.want.bulk {
				t.Errorf("handlePing(%v) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}

func TestHandleEcho(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want Value
	}{
		{
			name: "ECHO with one arg",
			args: []string{"hello"},
			want: Value{typ: "bulk", bulk: "hello"},
		},
		{
			name: "ECHO with no args",
			args: []string{},
			want: Value{typ: "error", str: "ERR wrong number of arguments, only one is expected"},
		},
		{
			name: "ECHO with two args",
			args: []string{"hello", "world"},
			want: Value{typ: "error", str: "ERR wrong number of arguments, only one is expected"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleEcho(tt.args)
			if got.typ != tt.want.typ || got.str != tt.want.str || got.bulk != tt.want.bulk {
				t.Errorf("handleEcho(%v) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}

func TestHandleSet(t *testing.T) {
	storage := NewStorage()

	tests := []struct {
		name string
		args []string
		want Value
	}{
		{
			name: "SET with two arguments",
			args: []string{"foo", "bar"},
			want: Value{typ: "string", str: "OK"},
		},
		{
			name: "SET with one argument",
			args: []string{"foo"},
			want: Value{typ: "error", str: "ERR wrong number of arguments, two are expected"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleSet(tt.args, storage)
			if got.typ != tt.want.typ || got.str != tt.want.str || got.bulk != tt.want.bulk {
				t.Errorf("handleSet(%v) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}

func TestHandleGet(t *testing.T) {
	storage := NewStorage()
	storage.Set("foo", "bar")

	tests := []struct {
		name string
		args []string
		want Value
	}{
		{
			name: "GET existing key",
			args: []string{"foo"},
			want: Value{typ: "bulk", bulk: "bar"},
		},
		{
			name: "GET non-existing key",
			args: []string{"baz"},
			want: Value{typ: "null"},
		},
		{
			name: "GET with no arguments",
			args: []string{},
			want: Value{typ: "error", str: "ERR wrong number of arguments, one is expected"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleGet(tt.args, storage)
			if got.typ != tt.want.typ || got.str != tt.want.str || got.bulk != tt.want.bulk {
				t.Errorf("handleGet(%v) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}

func TestHandleCommand(t *testing.T) {
	storage := NewStorage()

	tests := []struct {
		name string
		cmd  string
		args []string
		want Value
	}{
		{
			name: "PING command",
			cmd:  "PING",
			args: nil,
			want: Value{typ: "string", str: "PONG"},
		},
		{
			name: "ECHO command",
			cmd:  "ECHO",
			args: []string{"hi"},
			want: Value{typ: "bulk", bulk: "hi"},
		},
		{
			name: "Unknown command",
			cmd:  "FOO",
			args: []string{},
			want: Value{typ: "error", str: "ERR unknown command"},
		},
		{
			name: "Lowercase command still works",
			cmd:  "ping",
			args: nil,
			want: Value{typ: "string", str: "PONG"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleCommand(tt.cmd, tt.args, storage)
			if got.typ != tt.want.typ || got.str != tt.want.str || got.bulk != tt.want.bulk {
				t.Errorf("handleCommand(%q, %v) = %+v, want %+v", tt.cmd, tt.args, got, tt.want)
			}
		})
	}
}
