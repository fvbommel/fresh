package runner

import (
	"os"
	"path/filepath"

	"github.com/howeyc/fsnotify"
)

func watchFolder(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if isWatchedFile(ev.Name) {
					watcherLog("sending event %s", ev)
					startChannel <- ev.String()
				}
			case err := <-watcher.Error:
				watcherLog("error: %s", err)
			}
		}
	}()

	watcherLog("Watching %s", path)
	err = watcher.Watch(path)

	if err != nil {
		fatal(err)
	}
}

func watch() {
	watchRoot := watchRoot()
	filepath.Walk(watchRoot, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !isTmpDir(path) {
			// Skip hidden dirs other than "." and "..".
			base := filepath.Base(path)
			if len(base) > 1 && base[0] == '.' && base != ".." {
				return filepath.SkipDir
			}

			watchFolder(path)
		}

		return err
	})
}
