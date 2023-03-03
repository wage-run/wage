package wage_test

import (
	"os/exec"
	"testing"
	"time"

	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	"github.com/wage-run/wage"
)

func TestLoad(t *testing.T) {
	// resetTestdata()
	w := wage.NewWage("testdata")
	try.To(w.Start())
	_ = try.To1(w.Load("."))
	cmd := exec.Command("cp", "mod-wrong.go.tmpl", "mod.go")
	cmd.Dir = w.Root
	try.To(cmd.Run())
	time.Sleep(1 * time.Second)
	_, err := w.Load(".")
	assert.NotEqual(&err, nil)
}
