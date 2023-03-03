package wage

import (
	"testing"
	"time"

	"github.com/lainio/err2/try"
)

func TestCompile(t *testing.T) {
	resetTestdata()
	defer resetTestdata()
	w := NewWage("testdata")
	try.To(w.Start())
	p := w.getPkg(".")
	try.To1(w.Compile(p))
	time.Sleep(time.Second)
}
