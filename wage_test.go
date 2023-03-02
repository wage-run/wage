package wage

import (
	"testing"

	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
)

func TestFindPkg(t *testing.T) {
	w := NewWage("testdata")
	p := try.To1(w.FindPkg("foo/bar"))
	assert.Equal(p, "github.com/wage-run/wage/testdata/foo/bar")
}
