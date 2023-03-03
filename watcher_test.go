package wage

import (
	"os/exec"
	"testing"
	"time"

	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
)

func TestFsWatch(t *testing.T) {
	resetTestdata()
	defer resetTestdata()
	w := NewWage("testdata")
	w.Start()
	pkg := w.getPkg(".")
	dll1 := try.To1(w.Compile(pkg))
	dll2 := try.To1(w.Compile(pkg))
	assert.Equal(dll2, dll1)

	try.To(exec.Command("bash", "-c", "echo '//test' >> testdata/mod.go").Run())
	time.Sleep(time.Second)
	dll3 := try.To1(w.Compile(pkg))
	assert.NotEqual(dll3, dll1)
}

func resetTestdata() {
	cmd := exec.Command("git", "checkout", "testdata")
	try.To(cmd.Run())
}
