package sm

import (
	"io"
	"log"
	"os"
	"sync"
)

// AOF implementation

type AOFBuffer struct {
	Buffer   []byte
	Len      int
	Offset   int64
	FileName string
	File     *os.File
	Mutex    sync.Mutex
}

func NewAOFBuffer(filename string) *AOFBuffer {
	aof := &AOFBuffer{FileName: filename}
	aof.Offset = 0
	return aof
}

func (aof *AOFBuffer) Init() {
	var err error
	aof.File, err = os.OpenFile(aof.FileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	aof.Offset, err = aof.File.Seek(0, io.SeekEnd)
	if err != nil {
		log.Fatal(err.Error())
	}

}

func (aof *AOFBuffer) Write(cmd []byte) {
	aof.Mutex.Lock()
	aof.Buffer = append(aof.Buffer, cmd...)
	aof.Mutex.Unlock()

	if len(aof.Buffer) > 4194304 {
		aof.Sync()
	}
}

func (aof *AOFBuffer) Sync() {
	aof.Mutex.Lock()
	defer aof.Mutex.Unlock()

	if len(aof.Buffer) == 0 {
		return
	}

	n, err := aof.File.Write(aof.Buffer)
	if err != nil {
		return
	}

	aof.Buffer = aof.Buffer[n:]
	aof.Offset += int64(n)

	err = aof.File.Sync()
	if err != nil {
		return
	}
}

func (aof *AOFBuffer) Close() {
	for len(aof.Buffer) != 0 {
		aof.Sync()
	}
}
