package models

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

var Log *Logger

func init() {
	// 确保logs目录存在
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// 创建日志文件
	infoFile, err := os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open info log file:", err)
	}

	warningFile, err := os.OpenFile("logs/warning.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open warning log file:", err)
	}

	errorFile, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	// 同时写入文件和控制台
	infoMultiWriter := io.MultiWriter(os.Stdout, infoFile)
	warningMultiWriter := io.MultiWriter(os.Stdout, warningFile)
	errorMultiWriter := io.MultiWriter(os.Stderr, errorFile)

	Log = &Logger{
		infoLogger:    log.New(infoMultiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warningLogger: log.New(warningMultiWriter, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger:   log.New(errorMultiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.infoLogger.Println(v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.warningLogger.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.errorLogger.Println(v...)
}
