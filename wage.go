package wage

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"sync"
	"time"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/radovskyb/watcher"
)

type Wage struct {
	Root string

	tmpdir string
	mod    GoModInfo
	fsw    *watcher.Watcher

	rwl  *sync.RWMutex
	pkgs map[string]*pkg
}

type pkg struct {
	path     string
	compiled time.Time
	changed  time.Time
	locker   sync.Locker
}

func NewWage(root string) *Wage {
	if !filepath.IsAbs(root) {
		pwd := try.To1(os.Getwd())
		root = filepath.Join(pwd, root)
	}
	w := &Wage{
		Root: root,

		rwl:  &sync.RWMutex{},
		pkgs: map[string]*pkg{},
	}
	w.mod = try.To1(CollectGoModInfo(w.Root))
	w.tmpdir = filepath.Join(w.Root, ".wage-tmp")
	return w
}

func (w *Wage) Start() (err error) {
	defer err2.Handle(&err)
	try.To(os.MkdirAll(w.tmpdir, os.ModePerm))
	try.To(w.WatchFs())
	return
}

func (w *Wage) Close() (err error) {
	defer err2.Handle(&err)
	if w.fsw != nil {
		w.fsw.Close()
	}
	return
}

var ErrPkgNoExport = fmt.Errorf("pkg no Export func")

func (w *Wage) Load(path string) (m Module, err error) {
	defer err2.Handle(&err)
	pkg := try.To1(w.FindPkg(path))
	dll := try.To1(w.Compile(pkg))
	p := try.To1(plugin.Open(dll))
	export, ok := try.To1(p.Lookup("Export")).(Export)
	if !ok {
		err = fmt.Errorf("%w. %s", ErrPkgNoExport, pkg)
		return
	}
	m = export()
	return
}

func (w *Wage) FindPkg(path string) (string, error) {
	d := filepath.Join(w.Root, path)
	finfo, err := os.Stat(d)
	if err != nil {
		return "", err
	}
	if !finfo.IsDir() {
		d = filepath.Dir(d)
	}
	d = strings.TrimPrefix(d, w.mod.RootDir)
	d = filepath.Join(w.mod.ModulePath, d)
	return d, nil
}
