package engine

import (
	"bufio"
	"fmt"
	"os"
	"rodis/internal/protocol/resp"
	"sync"
	"time"
)

type Aof struct {
	File *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		File: f,
		rd:   bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk every 1 second
	go func() {
		for {
			aof.mu.Lock()

			aof.File.Sync()

			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.File.Close()
}

func (aof *Aof) Write(payload resp.Payload, rp resp.Resp) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	// _, err := aof.File.Write(rp.Marshal(payload))
	err := rp.Marshal(payload)
	if err != nil {
		return err
	}

	rp.FlushWriter()
	fmt.Println("writed request into aof file!")

	return nil
}

// func (aof *Aof) Read(callback func(payload resp.Payload)) error {
// 	aof.mu.Lock()
// 	defer aof.mu.Unlock()

// 	resp := resp.NewResp(aof.File)

// 	for {
// 		value, err := resp.Read()
// 		if err == nil {
// 			callback(value)
// 		}
// 		if err == io.EOF {
// 			break
// 		}
// 		return err
// 	}

// 	return nil
// }
