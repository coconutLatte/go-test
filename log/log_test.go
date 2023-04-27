package log

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	err := InitLogger(
		WithLevel(InfoLevel),
		WithUseConsole(true),
		WithPath("./test.log"))
	if err != nil {
		fmt.Println(err)
		return
	}

	Info("info")
	Debug("debug")
	Error("error")
}
