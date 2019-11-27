package sm

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type AOF struct {
	Buffer        []byte
	Mutex         sync.RWMutex
	SyncOffset    int32
	CurrentOffset int32
	File          *os.File
}

func LogIt(msg string) {
	log.Println(msg)
}

func ConvertInsert(name string, key string, value string) []byte {
	params := []string{
		"*4",
		name,
		key,
		key,
		value,
	}
	cmd := strings.Join(params, "\r\n")
	return []byte(cmd)
}

func NewAOF(filename string) *AOF {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)

	if err != nil {
		log.Fatal(err.Error())
	}

	aof := &AOF{}
	aof.File = file
	return aof
}

func (aof *AOF) Feed(cmd []byte) {
	aof.Mutex.Lock()
	aof.Buffer = append(aof.Buffer, cmd...)
	aof.CurrentOffset += int32(len(cmd))
	aof.Mutex.Unlock()
}

// Write buffer to disk
func (aof *AOF) Flush() {
	aof.Mutex.RLock()
	n, err := aof.File.Write(aof.Buffer)
	aof.Mutex.RUnlock()

	if err != nil {
		// log it
		LogIt(err.Error())
		return
	}

	aof.Mutex.Lock()
	aof.Buffer = aof.Buffer[n:]
	aof.SyncOffset = int32(n)
	aof.Mutex.Unlock()
}

func (aof *AOF) Sync() {
	err := aof.File.Sync()
	if err != nil {
		//log it
		LogIt(err.Error())
	}
}

func (aof *AOF) Close() {
	fmt.Println("close")
	fmt.Println(string(aof.Buffer))
	aof.Flush()
	aof.Sync()
	err := aof.File.Close()
	if err != nil {
		//log it
		LogIt(err.Error())
	}
}
