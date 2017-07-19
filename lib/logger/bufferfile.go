package logger

import (
	"bufio"
	"os"
)

type BufferFile struct {
	File    *os.File
	BufFile *bufio.Writer
}

func NewBufferFile(logfile string) (m *BufferFile, err error) {
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_SYNC, 0755)
	if err != nil {
		return
	}
	if file == nil {
		return
	}
	m = &BufferFile{file, bufio.NewWriter(file)}
	return
}

func (m *BufferFile) Close() {
	if m == nil {
		return
	}
	if m.File != nil {
		m.Flush()
		m.File.Close()
	}
}
func (m *BufferFile) Flush() {
	if m == nil {
		return
	}
	m.BufFile.Flush()
}

func (m *BufferFile) Write(p []byte) (n int, err error) {

	if m == nil {
		return
	}
	if p == nil || len(p) == 0 {
		return
	}
	return m.BufFile.Write(p)
}
