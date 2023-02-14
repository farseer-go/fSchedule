package test

import (
	"github.com/farseer-go/fSchedule"
	"github.com/farseer-go/fs"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	fs.Initialize[fSchedule.Module]("test fSchedule")
	defer fSchedule.Module{}.Shutdown()

	time.Sleep(3 * time.Second)
}
