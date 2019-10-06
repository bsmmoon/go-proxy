package logger

import (
	"bytes"
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/fatih/color"
)

var cmd string

func SetCmd(name string) {
	cmd = name
}

// INFO : low priority, no risk
func INFO(message string, a ...interface{}) {
	color.New(color.FgCyan).Println(fmt.Sprintf("[%v][INFO]", cmd), fmt.Sprintf(message, a...))
}

// WARNING : low priority, low risk
func WARNING(message string, a ...interface{}) error {
	fullMessage := fmt.Sprintf("[%v][WARNING] %v", cmd, fmt.Sprintf(message, a...))
	color.New(color.FgYellow).Println(fullMessage)
	return fmt.Errorf(fullMessage)
}

// ERROR : high priority, low rist
func ERROR(message string, a ...interface{}) error {
	fullMessage := fmt.Sprintf("[%v][ERROR] %v\n%v", cmd, fmt.Sprintf(message, a...), stack())
	color.New(color.FgRed).Println(fullMessage)
	return fmt.Errorf(fullMessage)
}

// FATAL : high priority, high risk. terminates the goroutine
func FATAL(message string, a ...interface{}) {
	color.New(color.FgRed).Println(fmt.Sprintf("[%v][FATAL]", cmd), fmt.Sprintf(message, a...), stack())
	log.Fatalln("FATAL!")
}

func stack() string {
	var buffer bytes.Buffer
	stack := string(debug.Stack())
	for _, stack := range strings.Split(stack, "\n") {
		if !strings.Contains(stack, "	/") {
			continue
		}
		if strings.Contains(stack, "stack.go") || strings.Contains(stack, "logger.go") {
			continue
		}
		buffer.WriteString(fmt.Sprintf("\n%v", stack))
	}
	return buffer.String()
}
