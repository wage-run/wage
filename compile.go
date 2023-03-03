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

	pkg.locker.Lock()
	defer pkg.locker.Unlock()

	if pkg.compiled.Sub(pkg.changed) >= 0 {
		return
	}
	pkg.compiled = time.Now()

	oldFilesMatch := strings.ReplaceAll(pkg.path, "/", "-") + "-*"
	oldFilesMatch = filepath.Join(w.tmpdir, oldFilesMatch)
	oldFiles := try.To1(filepath.Glob(oldFilesMatch))
	go func() {
		for _, f := range oldFiles {
			go os.Remove(f)
		}
	}()

	f := try.To1(os.OpenFile(gofile, os.O_CREATE|os.O_WRONLY, os.ModePerm))
	defer f.Close()
	mainTpl.Execute(f, map[string]any{"pkg": pkg.path})

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
