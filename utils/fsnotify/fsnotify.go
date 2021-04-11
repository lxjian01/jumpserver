package fsnotify

import (
	"github.com/howeyc/fsnotify"
	"sync"
)

var (
	once        sync.Once
	confOnce    sync.Once
	watcher     *fsnotify.Watcher
	confWatcher *fsnotify.Watcher
)

func GetWatcherOnce() (*fsnotify.Watcher, error) {
	var err error
	once.Do(func() {
		watcher, err = fsnotify.NewWatcher()
	})
	return watcher, err
}

func GetConfWatcherOnce() (*fsnotify.Watcher, error) {
	var err error
	confOnce.Do(func() {
		confWatcher, err = fsnotify.NewWatcher()
	})
	return confWatcher, err
}
