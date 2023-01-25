package test

import (
	"github.com/farseer-go/fSchedule"
	"github.com/farseer-go/fs"
	"testing"
)

func TestRun(t *testing.T) {
	fs.Initialize[fSchedule.Module]("test fSchedule")
	fSchedule.Module{}.Shutdown()
}
