package wage

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/radovskyb/watcher"
)

func (w *Wage) WatchFs() (err error) {
	defer err2.Handle(&err)
	w.fsw = watcher.New()
	fw := w.fsw
	try.To(fw.AddRecursive(w.Root))
	go func() {
		for {
			select {
			case <-fw.Closed:
				return
			case <-fw.Error:
				// fmt.Println("wa")
			case e := <-fw.Event:
				if strings.HasPrefix(e.Path, w.tmpdir) {
					continue
				}
				if e.FileInfo.IsDir() {
					go w.handleDirChange(e)
				} else {
					go w.handleFileChange(e)
				}
			}
		}
	}()
	go fw.Start(200 * time.Microsecond)
	return
}

func (w *Wage) handleDirChange(e watcher.Event) (err error) {
	w.pkgsL.Lock()
	defer w.pkgsL.Unlock()
	var key string
	switch e.Op {
	case watcher.Move:
		key = w.FindPkg(e.OldPath)
		delete(w.pkgs, key)
	case watcher.Remove:
		key = w.FindPkg(e.Path)
		delete(w.pkgs, key)
	}
	return
}

func (w *Wage) handleFileChange(e watcher.Event) (err error) {

	if !isGoFile(e.Path) {
		return
	}

	dir := filepath.Dir(e.Path)
	pkg := w.getPkg(w.FindPkg(dir))

	if pkg == nil {
		return
	}

	pkg.locker.Lock()
	defer pkg.locker.Unlock()

	switch e.Op {
	case watcher.Chmod:
		return
	default:
		pkg.changed = time.Now()
	}
	return
}

func isGoFile(f string) bool {
	return strings.HasSuffix(f, ".go")
}
