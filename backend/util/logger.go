package util

import (
	"log"
	"os"
	"path/filepath"
)

var logFile *os.File

func InitLogger(logDir string) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("failed to create log dir: %v", err)
	}

	logPath := filepath.Join(logDir, "server.log")

	// 启动时清空日志文件
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	logFile = f
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func GetLogFilePath() string {
	if logFile == nil {
		return ""
	}
	return logFile.Name()
}
