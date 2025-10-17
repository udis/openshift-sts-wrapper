package logger

import (
	"fmt"
	"io"
	"os"
)

type Level int

const (
	LevelQuiet Level = iota
	LevelNormal
	LevelVerbose
)

type Logger struct {
	level  Level
	writer io.Writer
}

func New(level Level, writer io.Writer) *Logger {
	if writer == nil {
		writer = os.Stdout
	}
	return &Logger{
		level:  level,
		writer: writer,
	}
}

func (l *Logger) Info(msg string) {
	if l.level >= LevelNormal {
		fmt.Fprintln(l.writer, msg)
	}
}

func (l *Logger) Debug(msg string) {
	if l.level >= LevelVerbose {
		fmt.Fprintln(l.writer, msg)
	}
}

func (l *Logger) Error(msg string) {
	fmt.Fprintln(l.writer, msg)
}

func (l *Logger) StartStep(name string) {
	if l.level >= LevelNormal {
		fmt.Fprintf(l.writer, "⏳ %s...\n", name)
	}
}

func (l *Logger) CompleteStep(name string) {
	if l.level >= LevelNormal {
		fmt.Fprintf(l.writer, "✓ %s\n", name)
	}
}

func (l *Logger) FailStep(name string) {
	fmt.Fprintf(l.writer, "✗ %s\n", name)
}
