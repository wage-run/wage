package wage

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type GoModInfo struct {
	RootDir    string // go.mod dir
	ModulePath string // go.mod module path
}

func CollectGoModInfo(workdir string) (info GoModInfo, err error) {
	defer err2.Handle(&err)

	{
		var b bytes.Buffer
		cmd := exec.Command("go", "env", "GOMOD")
		cmd.Dir = workdir
		cmd.Stdout = &b
		try.To(cmd.Run())
		info.RootDir = filepath.Dir(b.String())
	}

	{
		var b bytes.Buffer
		cmd := exec.Command("go", "mod", "edit", "-json")
		cmd.Dir = workdir
		cmd.Stdout = &b
		try.To(cmd.Run())
		var output struct {
			Module struct {
				Path string
			}
		}
		json.Unmarshal(b.Bytes(), &output)
		info.ModulePath = output.Module.Path
	}
	return
}
