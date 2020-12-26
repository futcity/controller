///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package utils

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Log types
const (
	LogInfo = iota
	LogWarn
	LogError
)

// Log Logger struct
type Log struct {
	path string
}

// NewLog make new struct
func NewLog() *Log {
	return &Log{}
}

// SetPath Set log path
func (l *Log) SetPath(path string) {
	l.path = path
}

// LogMessage log save message
func (l *Log) LogMessage(module string, msg string, error string, logType int) error {
	var outMsg string
	var typ string
	dt := time.Now()

	// Configuring log mod
	var log, err = os.OpenFile(l.path+dt.Format("20060201")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println("Fail to open log file!")
		return err
	}

	switch logType {
	case LogInfo:
		typ = "INFO"
		break

	case LogWarn:
		typ = "WARN"
		break

	case LogError:
		typ = "ERROR"
		break
	}

	outMsg = "[" + dt.Format("15:04:05") + "][" + module + "][" + typ + "] "
	outMsg += msg

	if logType == LogError {
		outMsg += ": " + error
	}

	fmt.Println(outMsg)

	var _, errWr = log.WriteString(outMsg + "\n")
	if errWr != nil {
		fmt.Println("Fail to write to logfile!")
		log.Close()
		return errWr
	}

	log.Close()
	return nil
}

// Info Information log message
func (l *Log) Info(module string, msg string) error {
	return l.LogMessage(module, msg, "", LogInfo)
}

// Error Error log message
func (l *Log) Error(module string, msg string, err string) error {
	return l.LogMessage(module, msg, err, LogError)
}

// TodayMessages get today messages
func (l *Log) TodayMessages() ([]string, error) {
	var data []string
	dt := time.Now()

	file, err := os.Open(l.path + dt.Format("20060201") + ".log")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}
