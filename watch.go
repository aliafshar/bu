package bu

import(
  "golang.org/x/exp/inotify"
	"github.com/aliafshar/toylog"
)

type watcher struct {
  path string
  out chan bool
}


func (w *watcher) watch() error {
  watcher, err := inotify.NewWatcher()
  if err != nil {
    return err
  }
  err = watcher.Watch(w.path)
  if err != nil {
    return err
  }
  for {
    select {
    case ev := <-watcher.Event:
        w.out <- true
        toylog.Debugf("change detected %v", ev)
    case err := <-watcher.Error:
        w.out <- true
        toylog.Errorf("error watching %q, %q\n", w.path, err)
    }
  }
  return nil
}

