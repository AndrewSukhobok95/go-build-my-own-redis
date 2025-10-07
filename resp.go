package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

func (v Value) marshalString() []byte {
	bytes := []byte{'+'}
	bytes = append(bytes, v.str...)
	bytes = append(bytes, "\r\n"...)
	return bytes
}

func (v Value) marshalInteger() []byte {
	bytes := []byte{':'}
	bytes = append(bytes, strconv.Itoa(v.num)...)
	bytes = append(bytes, "\r\n"...)
	return bytes
}

func (v Value) marshalError() []byte {
	bytes := []byte{'-'}
	bytes = append(bytes, v.str...)
	bytes = append(bytes, "\r\n"...)
	return bytes
}

func (v Value) marshalBulk() []byte {
	bytes := []byte{'$'}
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, "\r\n"...)
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, "\r\n"...)
	return bytes
}

func (v Value) marshalArray() []byte {
	bytes := []byte{'*'}
	bytes = append(bytes, strconv.Itoa(len(v.array))...)
	bytes = append(bytes, "\r\n"...)
	for i := 0; i < len(v.array); i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}
	return bytes
}

func (v Value) Marshal() []byte {
	switch v.typ {
	case "bulk":
		return v.marshalBulk()
	case "array":
		return v.marshalArray()
	case "string":
		return v.marshalString()
	case "integer":
		return v.marshalInteger()
	case "error":
		return v.marshalError()
	default:
		return nil
	}
}

func parseCommand(v Value) (cmd string, args []string, err error) {
	if v.typ != "array" {
		return "", nil, fmt.Errorf("invalid RESP type: %s", v.typ)
	}
	if len(v.array) == 0 {
		return "", nil, fmt.Errorf("empty command")
	}
	cmd = v.array[0].bulk
	for i := 1; i < len(v.array); i++ {
		args = append(args, v.array[i].bulk)
	}
	return cmd, args, nil
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	var b byte
	for {
		b, err = r.reader.ReadByte()
		if err != nil {
			return line, n, err
		}
		n++
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
			line = line[:len(line)-2]
			return line, n, nil
		}
	}
}

func (r *Resp) readInteger() (x int64, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, n, err
	}
	s := string(line)
	x, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, n, err
	}
	return x, n, nil
}

func (r *Resp) readBulk() (Value, error) {
	x, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}
	if x < 0 {
		return Value{typ: "bulk", bulk: ""}, nil
	}
	s := make([]byte, x+2)
	_, err = io.ReadFull(r.reader, s)
	if err != nil {
		return Value{}, err
	}
	return Value{typ: "bulk", bulk: string(s[:len(s)-2])}, nil
}

func (r *Resp) readArray() (Value, error) {
	x, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}
	if x == -1 {
		return Value{typ: "array", array: nil}, nil
	}
	if x < -1 {
		return Value{}, fmt.Errorf("invalid array length: %d", x)
	}
	array := make([]Value, 0, x)
	for i := int64(0); i < x; i++ {
		v, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		array = append(array, v)
	}
	return Value{typ: "array", array: array}, nil
}

func (r *Resp) Read() (Value, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch b {
	case BULK:
		return r.readBulk()
	case ARRAY:
		return r.readArray()
	default:
		return Value{}, fmt.Errorf("unknown RESP type: %q", b)
	}
}

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	bytes := v.Marshal()
	_, err := w.writer.Write(bytes)
	return err
}
