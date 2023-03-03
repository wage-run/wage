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

func (w *Wage) Compile(pkg *pkg) (dll string, err error) {
	defer err2.Handle(&err)
	if pkg == nil {
		return "", os.ErrNotExist
	}

	var gofile string
	gofile = strings.ReplaceAll(pkg.path, "/", "-") + fmt.Sprintf("-%d.go", pkg.changed.Unix())
	gofile = filepath.Join(w.tmpdir, gofile)
	dll = strings.TrimSuffix(gofile, ".go") + ".so"

	if !needRecompile(pkg) {
		return
	}
	pkg.compiled = time.Now()

	f := try.To1(os.OpenFile(gofile, os.O_CREATE|os.O_WRONLY, os.ModePerm))
	try.To(mainTpl.Execute(f, map[string]any{"pkg": pkg.path}))
	f.Close()

	cmd := exec.Command("go", "build", "-buildmode=plugin", gofile)
	cmd.Stderr = os.Stderr
	cmd.Dir = w.tmpdir
	try.To(cmd.Run())

	return
}

func (w *Wage) getPkg(pkgPath string) (p *pkg) {
	pkgPath = w.FindPkg(pkgPath)
	w.pkgsL.RLock()
	p, ok := w.pkgs[pkgPath]
	w.pkgsL.RUnlock()
	if !ok {
		w.pkgsL.Lock()
		defer w.pkgsL.Unlock()

		p = &pkg{
			path:    pkgPath,
			changed: time.Now(),
			locker:  &sync.Mutex{},
		}

		w.pkgs[pkgPath] = p

		return
	}

	return
}

func needRecompile(pkg *pkg) bool {
	pkg.locker.Lock()
	defer pkg.locker.Unlock()
	t := pkg.compiled.Sub(pkg.changed)
	return t < 0
}
