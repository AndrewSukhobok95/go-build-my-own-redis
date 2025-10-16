package commands

import (
	"testing"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

func TestHandlePing(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want resp.Value
	}{
		{
			name: "PING with no args",
			args: nil,
			want: resp.NewStringValue("PONG"),
		},
		{
			name: "PING with one arg",
			args: []string{"hello"},
			want: resp.NewBulkValue("hello"),
		},
		{
			name: "PING with too many args",
			args: []string{"foo", "bar"},
			want: resp.NewErrorValue("ERR wrong number of arguments"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handlePing(tt.args)
			if got.Typ() != tt.want.Typ() || got.Str() != tt.want.Str() || got.Bulk() != tt.want.Bulk() {
				t.Errorf("handlePing(%v) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}

func TestHandleEcho(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want resp.Value
	}{
		{
			name: "ECHO with one arg",
			args: []string{"hello"},
			want: resp.NewBulkValue("hello"),
		},
		{
			name: "ECHO with no args",
			args: []string{},
			want: resp.NewErrorValue("ERR wrong number of arguments"),
		},
		{
			name: "ECHO with two args",
			args: []string{"hello", "world"},
			want: resp.NewErrorValue("ERR wrong number of arguments"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleEcho(tt.args)
			if got.Typ() != tt.want.Typ() || got.Str() != tt.want.Str() || got.Bulk() != tt.want.Bulk() {
				t.Errorf("handleEcho(%v) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}

func TestHandleSet(t *testing.T) {
	storage := storage.NewKV()

	tests := []struct {
		name string
		args []string
		want resp.Value
	}{
		{
			name: "SET with two arguments",
			args: []string{"foo", "bar"},
			want: resp.NewStringValue("OK"),
		},
		{
			name: "SET with one argument",
			args: []string{"foo"},
			want: resp.NewErrorValue("ERR wrong number of arguments"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleSet(tt.args, storage)
			if got.Typ() != tt.want.Typ() || got.Str() != tt.want.Str() || got.Bulk() != tt.want.Bulk() {
				t.Errorf("handleSet(%v) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}

func TestHandleGet(t *testing.T) {
	storage := storage.NewKV()
	storage.Set("foo", "bar")

	tests := []struct {
		name string
		args []string
		want resp.Value
	}{
		{
			name: "GET existing key",
			args: []string{"foo"},
			want: resp.NewBulkValue("bar"),
		},
		{
			name: "GET non-existing key",
			args: []string{"baz"},
			want: resp.NewNullValue(),
		},
		{
			name: "GET with no arguments",
			args: []string{},
			want: resp.NewErrorValue("ERR wrong number of arguments"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleGet(tt.args, storage)
			if got.Typ() != tt.want.Typ() || got.Str() != tt.want.Str() || got.Bulk() != tt.want.Bulk() {
				t.Errorf("handleGet(%v) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}

func TestHandleCommand(t *testing.T) {
	storage := storage.NewKV()

	tests := []struct {
		name string
		cmd  string
		args []string
		want resp.Value
	}{
		{
			name: "PING command",
			cmd:  "PING",
			args: nil,
			want: resp.NewStringValue("PONG"),
		},
		{
			name: "ECHO command",
			cmd:  "ECHO",
			args: []string{"hi"},
			want: resp.NewBulkValue("hi"),
		},
		{
			name: "Unknown command",
			cmd:  "FOO",
			args: []string{},
			want: resp.NewErrorValue("ERR unknown command"),
		},
		{
			name: "Lowercase command still works",
			cmd:  "ping",
			args: nil,
			want: resp.NewStringValue("PONG"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HandleCommand(tt.cmd, tt.args, storage)
			if got.Typ() != tt.want.Typ() || got.Str() != tt.want.Str() || got.Bulk() != tt.want.Bulk() {
				t.Errorf("handleCommand(%q, %v) = %+v, want %+v", tt.cmd, tt.args, got, tt.want)
			}
		})
	}
}
