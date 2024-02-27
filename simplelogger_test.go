package simplelogger

import (
	"log"
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	logfile := "./log_test.txt"

	logger := &Logger{MaxSize: 100}
	if err := logger.Open(logfile); err != nil {
		t.Fatal(err)
	}

	defer func() {
		os.Remove(logfile)
	}()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Println("hello world: ", i)
	}
}
