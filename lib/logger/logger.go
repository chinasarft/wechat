package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	FATAL = iota
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)

var levelStr = []string{
	"FATAL",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
	"TRACE",
}

type JsonLogger struct {
	mu        sync.Mutex
	callDepth int
	// loglevel is a interger type, so comparison between loglevels are supported
	level int
	// log writer will be sorted by priority
	logger *BufferFile

	path       string //log file's parent path
	size       int64  //current file size
	backupNum  int    //backup num
	rotateSize int64
}

type logFormat struct {
	Time  string `bson:"time"    json:"time"`
	Level string `bson:"level"   json:"level"`
	Flag  string `bson:"flag"    json:"flag"`
	File  string `bson:"file"    json:"file"`
	Line  int    `bson:"line"    json:"line"`
	Msg   string `bson:"message" json:"message"`
}

// -----------------------------------------------------------------------------------------------------------

func NewLogger(dir, localFile string) (jl *JsonLogger, err error) {
	var st os.FileInfo
	jl = &JsonLogger{}
	jl.callDepth = 3
	jl.level = INFO
	jl.backupNum = 3
	jl.rotateSize = 1024 * 1024 * 50

	if dir == "" {
		jl.path = os.Getenv("PWD")
	} else {
		jl.path = dir
	}

	jl.logger, err = NewBufferFile(path.Join(jl.path, localFile))
	st, err = jl.logger.File.Stat()
	if err != nil {
		return
	}
	jl.size = st.Size()

	return
}

func (l *JsonLogger) SetLogLevel(level int) {
	if level < 0 || level > TRACE {
		return
	}
	l.level = level
	return
}

func (l *JsonLogger) SetCallDepth(callDepth int) {
	if callDepth > 0 {
		l.callDepth = callDepth
	}
}

func (l *JsonLogger) SetRotateSize(len int64) {
	if len > 0 {
		l.rotateSize = len
	}
}

func (l *JsonLogger) SetBackupNum(num int) {
	if num > 0 {
		l.backupNum = num
	}
}

func (l *JsonLogger) Flush() {
	if l != nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Flush()
	}
}

func (l *JsonLogger) Close() {
	if l != nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Close()
	}
}

func (m *JsonLogger) rotate() (err error) {
	var i int
	var tarFile string
	var bakFile []string

	if m.size < m.rotateSize {
		return
	}
	fmt.Println("roate:", m.size, m.rotateSize)
	m.logger.Close()
	filename := m.logger.File.Name()

	for i = m.backupNum; i > 0; i-- {
		tarFile = filename + "." + strconv.Itoa(i)
		bakFile = append(bakFile, tarFile)
		_, err = os.Stat(tarFile)
		if i == m.backupNum && err == nil {
			if m.path == "" {
				os.Remove(tarFile)
			} else {

				os.Remove(path.Join(m.path, tarFile))
			}
			continue
		}
		if err == nil {
			err = os.Rename(tarFile, bakFile[m.backupNum-i-1])
			if err != nil {
				return
			}
		}
	}
	err = os.Rename(filename, bakFile[m.backupNum-1])
	if err != nil {
		return
	}
	m.logger, err = NewBufferFile(filename)
	m.size = 0
	return
}

func (l *JsonLogger) writeLog(level int, flag, msg string, v ...interface{}) (n int, err error) {

	// do not log message if loglevel does not match
	if level > l.level {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	err = l.rotate()
	if err != nil {
		return
	}

	// foramt output strings
	_, file, line, _ := runtime.Caller(l.callDepth)
	logmsg := msg
	if v != nil && len(v) > 0 {
		logmsg = fmt.Sprintf(msg, v...)
	}
	bs, err := json.Marshal(
		&logFormat{
			Time:  time.Now().Format("2006-01-02 15:04:05.000000000"),
			Level: levelStr[level],
			Flag:  flag,
			File:  path.Base(file),
			Line:  line,
			Msg:   logmsg,
		},
	)
	if err != nil {
		return
	}

	n, err = l.logger.Write(bs)
	if err == nil {
		l.logger.Write([]byte("\n"))
		l.size += int64(n + 1)
	}
	return
}

// -----------------------------------------------------------------------------------------------------------

func (l *JsonLogger) Fatal(level int, flag, msg string, v ...interface{}) {
	l.writeLog(FATAL, flag, msg, v...)
}

func (l *JsonLogger) Error(level int, flag, msg string, v ...interface{}) {
	l.writeLog(ERROR, flag, msg, v...)
}

func (l *JsonLogger) Warn(level int, flag, msg string, v ...interface{}) {
	l.writeLog(WARN, flag, msg, v...)
}

func (l *JsonLogger) Info(level int, flag, msg string, v ...interface{}) {
	l.writeLog(INFO, flag, msg, v...)
}

func (l *JsonLogger) Debug(level int, flag, msg string, v ...interface{}) {
	l.writeLog(DEBUG, flag, msg, v...)
}

func (l *JsonLogger) Trace(level int, flag, msg string, v ...interface{}) {
	l.writeLog(TRACE, flag, msg, v...)
}
