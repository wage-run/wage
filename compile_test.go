package wage

import (
	"testing"

	"github.com/lainio/err2/try"
)

func TestCompile(t *testing.T) {
	w := NewWage("testdata")
	try.To(w.Start())
	p := try.To1(w.FindPkg("."))
	try.To1(w.Compile(p))
}
