package logger

import (
	"strconv"
	"testing"
)

func TestDefaultRotate(t *testing.T) {

	err := Start("", "tmp.log")
	if err != nil {
		t.Fatal("start logger fail:", err)
	}
	SetRotateSize(8096)
	for i := 0; i < 1000; i++ {
		Info(strconv.Itoa(i), "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	}
	Close()
}
