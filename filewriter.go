package crawlfrontier

import (
	"log"
	"os"
	"fmt"
	"time"
)

type Logger struct {
	file string
}

func NewLogger() *Logger {
	path := "~/"
	dir, err := os.Getwd()
	if err == nil {
		path = dir
	} else {
		fmt.Println("err," + err.Error())
	}
	
	dt := time.Now()
	filename := "frontier_" + dt.Format("2006-01-02") + ".log"
	filepath := path + "/" + filename
	fmt.Println("Log file path: ",filepath)
	logger := Logger{
		file: filepath,
	}
	return &logger
}

func (logger *Logger) writeLog(content ...string) error {
	logger.checkLogName()  // if date change, update log filename
	f, err := os.OpenFile(logger.file, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	for _, text := range content {
		log.Println(text)
	}
	return nil
}

func (logger *Logger) checkLogName() {
	dt := time.Now()
	filename := "frontier_" + dt.Format("2006-01-02") + ".log"
	if logger.file != filename {
		logger.file = filename
	}
}