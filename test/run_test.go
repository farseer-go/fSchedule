package test

import (
	"github.com/farseer-go/fSchedule"
	"github.com/farseer-go/fs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun(t *testing.T) {
	assert.NotPanics(t, func() {
		fSchedule.Module{}.Shutdown()
	})
	fs.Initialize[fSchedule.Module]("test fSchedule")

	client := fSchedule.GetClient()
	client.LogoutClient()

	//time.Sleep(1 * time.Hour)
}
