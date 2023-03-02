package wage

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

//go:embed compile-main.tmpl
var mainTplStr string

var mainTpl = template.Must(template.New("tmp-main").Parse(mainTplStr))

func (w *Wage) Compile(pkgPath string) (dll string, err error) {
	defer err2.Handle(&err)
	pkg := w.getPkg(pkgPath)
	if pkg == nil {
		return "", os.ErrNotExist
	}
	defer func() { dll = strings.TrimSuffix(dll, ".go") + ".so" }()
	dll = strings.ReplaceAll(pkg.path, "/", "-") + fmt.Sprintf("-%d.go", pkg.changed.Unix())
	dll = filepath.Join(w.tmpdir, dll)
	if !pkg.compiled.IsZero() {
		return
	}

	pkg.locker.Lock()
	defer pkg.locker.Unlock()

	f := try.To1(os.OpenFile(dll, os.O_CREATE|os.O_WRONLY, os.ModePerm))
	defer f.Close()
	mainTpl.Execute(f, map[string]any{"pkg": pkgPath})

	cmd := exec.Command("go", "build", "-buildmode=plugin", dll)
	cmd.Stderr = os.Stderr
	cmd.Dir = w.tmpdir
	try.To(cmd.Run())

	return
}

func (w *Wage) getPkg(pkgPath string) (p *pkg) {
	w.rwl.Lock()
	defer w.rwl.Unlock()

	p, ok := w.pkgs[pkgPath]
	if !ok {
		pkgFsPath := w.mod.Dir(pkgPath)
		stat, err := os.Stat(pkgFsPath)
		if err != nil {
			return nil
		}
		p = &pkg{
			path:    pkgPath,
			changed: stat.ModTime(),
			locker:  &sync.Mutex{},
		}
		return
	}

	if p.changed.Sub(p.compiled) >= 0 {
		p.compiled = time.Time{}
	}

	return
}
