package simplelogger

import (
	"os"
	"sync"
)

// TODO 目前是单个日志文件满了后直接清空，后续支持多文件
type Logger struct {
	MaxSize int64 // 单个日志文件的大小限制，单位：MB

	filename string
	size     int64
	file     *os.File
	mutex    sync.Mutex
}

func (l *Logger) Open(filename string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.open(filename)
}

func (l *Logger) open(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	info, err := os.Stat(filename)
	if err != nil {
		return err
	}

	l.filename = filename
	l.size = info.Size()
	l.file = file

	return nil
}

func (l *Logger) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.file == nil {
		return nil
	}

	err := l.file.Close()
	l.file = nil

	return err
}

func (l *Logger) Write(p []byte) (int, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.file == nil {
		return len(p), nil
	}

	if l.size > l.MaxSize*1024*1024 {
		l.file.Truncate(0)
		l.file.Close()
		l.open(l.filename)
	}

	n, err := l.file.Write(p)
	l.size += int64(n)

	return n, err
}
