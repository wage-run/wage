package wage

import (
	"os"
	"testing"

	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
)

func TestCollectGoModInfo(t *testing.T) {
	workdir := try.To1(os.Getwd())
	minfo := try.To1(CollectGoModInfo(workdir))
	assert.Equal(minfo.ModulePath, "github.com/wage-run/wage")
	assert.Equal(minfo.RootDir, workdir)
}
