package persistence

import (
	"bytes"
	"io"
	"os"
	"sync"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/engine"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

type AOF struct {
	mu     sync.RWMutex
	file   *os.File
	writer *resp.Writer
	buf    *bytes.Buffer
	path   string
}

func NewAOF(path string) (*AOF, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	w := resp.NewWriter(f)
	return &AOF{
		file:   f,
		writer: w,
		buf:    &bytes.Buffer{},
		path:   path,
	}, nil
}

func (aof *AOF) Append(cmdName string, args []string) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	array := make([]resp.Value, len(args)+1)
	array[0] = resp.NewBulkValue(cmdName)
	for i, arg := range args {
		array[i+1] = resp.NewBulkValue(arg)
	}
	respArray := resp.NewArrayValue(array)

	_, err := aof.buf.Write(respArray.Marshal())
	return err
}

func (aof *AOF) Flush() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	if aof.buf.Len() == 0 {
		return nil
	}
	_, err := aof.file.Write(aof.buf.Bytes())
	if err != nil {
		return err
	}
	aof.buf.Reset()
	return aof.file.Sync()
}

func (aof *AOF) Close() error {
	errFlush := aof.Flush()
	errClose := aof.file.Close()
	if errFlush != nil {
		return errFlush
	}
	if errClose != nil {
		return errClose
	}
	return nil
}

func (aof *AOF) Load(storage *storage.KV) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	f, err := os.OpenFile(aof.path, os.O_RDONLY, 0644)
	switch {
	case os.IsNotExist(err):
		return nil
	case err != nil:
		return err
	}

	ctx := engine.NewCommandContext(storage)
	ctx.StartReplay()

	respReader := resp.NewReader(f)
	for {
		v, err := respReader.Read()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		}
		cmdName, args, err := resp.ParseCommand(v)
		if err != nil {
			return err
		}
		engine.DispatchCommand(ctx, cmdName, args)
	}
}
