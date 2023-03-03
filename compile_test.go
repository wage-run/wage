package wage

import (
	"testing"
	"time"

	"github.com/lainio/err2/try"
)

func TestCompile(t *testing.T) {
	w := NewWage("testdata")
	try.To(w.Start())
	p := w.getPkg(w.FindPkg("."))
	try.To1(w.Compile(p))
	time.Sleep(time.Second)
}
