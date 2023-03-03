package wage

import (
	"testing"

	"github.com/lainio/err2/assert"
	// "github.com/wage-run/wage"
)

func TestFindPkg(t *testing.T) {
	w := NewWage("testdata")
	p := w.FindPkg("foo/bar")
	assert.Equal(p, "github.com/wage-run/wage/testdata/foo/bar")
}
