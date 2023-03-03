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

	pkgsL *sync.RWMutex
	pkgs  map[string]*pkg
}

type pkg struct {
	path     string
	compiled time.Time
	changed  time.Time
	locker   sync.Locker
}

func NewWage(root string) *Wage {
	root = try.To1(filepath.Abs(root))

	w := &Wage{
		Root: root,

		pkgsL: &sync.RWMutex{},
		pkgs:  map[string]*pkg{},
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
	pkg := w.getPkg(path)
	dll := try.To1(w.Compile(pkg))
	p := try.To1(plugin.Open(dll))
	export, ok := try.To1(p.Lookup("Export")).(Export)
	if !ok {
		err = fmt.Errorf("%w. %s", ErrPkgNoExport, pkg.path)
		return
	}
	m = export()
	return
}

func (w *Wage) FindPkg(dir string) string {
	if strings.HasPrefix(dir, w.mod.ModulePath) {
		return dir
	}
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(w.Root, dir)
	}
	dir = strings.TrimPrefix(dir, w.mod.RootDir)
	dir = filepath.Join(w.mod.ModulePath, dir)
	return dir
}
