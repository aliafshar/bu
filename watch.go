package bu

import (
	"github.com/aliafshar/toylog"
	"golang.org/x/exp/inotify"
)

type watcher struct {
	paths []string
	out   chan bool
}

func (w *watcher) watch() error {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		return err
	}
	for _, p := range w.paths {
		err = watcher.Watch(p)
		if err != nil {
			return err
		}
	}
	for {
		select {
		case ev := <-watcher.Event:
			if ev.Mask&inotify.IN_ATTRIB > 0 {
				toylog.Debugf("change detected %v", ev)
				w.out <- true
			}
		case err := <-watcher.Error:
			w.out <- true
			toylog.Errorf("error watching %q, %q\n", w.paths, err)
		}
	}
	return nil
}
