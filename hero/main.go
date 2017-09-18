package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/howeyc/fsnotify"
	"github.com/shiyanhui/hero"
)

var (
	watch                 bool
	source, dest, pkgName string
)

func init() {
	flag.StringVar(
		&source,
		"source",
		"./",
		"the html template file or dir",
	)
	flag.StringVar(
		&dest,
		"dest",
		"",
		"generated golang files dir, it will be the same with source if not set",
	)
	flag.StringVar(
		&pkgName,
		"pkgname",
		"template",
		"the generated template package name, default is `template`",
	)
	flag.BoolVar(
		&watch,
		"watch",
		false,
		"whether automatically compile when the source files change",
	)
}

func watchFile(watcher *fsnotify.Watcher, path string) {
	if err := watcher.Watch(path); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if dest == "" {
		dest = source
	}

	hero.Generate(source, dest, pkgName)

	if !watch {
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for ev := range watcher.Event {
			if ev.IsDelete() || ev.IsModify() || ev.IsRename() {
				hero.Generate(source, dest, pkgName)
			}
		}
	}()

	watchFile(watcher, source)

	stat, _ := os.Stat(source)
	if stat.IsDir() {
		filepath.Walk(source, func(
			path string, _ os.FileInfo, err error) error {

			stat, _ := os.Stat(path)
			if stat.IsDir() {
				watchFile(watcher, path)
			}
			return nil
		})
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done

	watcher.Close()
}
