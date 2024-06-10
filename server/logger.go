package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	InfoLog    *log.Logger
	WarningLog *log.Logger
	ErrorLog   *log.Logger
	FatalLog   *log.Logger
)

func sendLog(message string) {
	addr := fmt.Sprintf(":%d", LOGGER_PORT)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("Failed to connect to logger server:", err)
		os.Exit(1)
		return
	}
	defer conn.Close()
	if _, err := conn.Write([]byte(message)); err != nil {
		log.Println("Failed to send log message:", err)
	}
}

func init() {
	mw := io.MultiWriter(os.Stdout)
	InfoLog = log.New(mw, "[info ]: ", log.Ldate|log.Ltime)
	WarningLog = log.New(mw, "[warn ]: ", log.Ldate|log.Ltime)
	ErrorLog = log.New(mw, "[error]: ", log.Ldate|log.Ltime)
	FatalLog = log.New(mw, "[fatal]: ", log.Ldate|log.Ltime)

}

func Info(format string, v ...interface{}) {
	payload := fmt.Sprintf(format, v...)
	InfoLog.Printf(payload)
	sendLog(payload)
}

func Warning(format string, v ...interface{}) {
	payload := fmt.Sprintf(format, v...)
	WarningLog.Printf(payload)
	sendLog(payload)
}

func Error(format string, v ...interface{}) {
	payload := fmt.Sprintf(format, v...)
	ErrorLog.Printf(payload)
	sendLog(payload)
}

func Fatal(format string, v ...interface{}) {
	payload := fmt.Sprintf(format, v...)
	FatalLog.Printf(payload)
	sendLog(payload)
	os.Exit(1)
}
