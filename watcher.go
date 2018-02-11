package inotify

import (
	"errors"
	"strings"
	"syscall"
	"unsafe"
)

type Watcher struct {
	fd int

	pathToWdMap map[string]uint32
	wdToPathMap map[uint32]string
}

func NewWatcher() (*Watcher, error) {
	fd, err := syscall.InotifyInit()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		fd: fd,

		pathToWdMap: make(map[string]uint32),
		wdToPathMap: make(map[uint32]string),
	}, nil
}

func (w *Watcher) AddWatch(path string, mask uint32) error {
	path = strings.TrimRight(path, "/")
	_, ok := w.pathToWdMap[path]
	if ok {
		return nil
	}

	wd, err := syscall.InotifyAddWatch(w.fd, path, mask)
	if err != nil {
		return err
	}

	uwd := uint32(wd)
	w.pathToWdMap[path] = uwd
	w.wdToPathMap[uwd] = path

	return nil
}

func (w *Watcher) RmWatch(path string) {
	path = strings.TrimRight(path, "/")
	wd, ok := w.pathToWdMap[path]
	if !ok {
		return
	}

	syscall.InotifyRmWatch(w.fd, wd)
	delete(w.pathToWdMap, path)
	delete(w.wdToPathMap, wd)
}

func (w *Watcher) ReadEvents() ([]*Event, error) {
	buf := make([]byte, syscall.SizeofInotifyEvent*4096)

	n, err := syscall.Read(w.fd, buf)
	if n == 0 {
		return nil, errors.New("Read 0 byte error")
	}
	if err != nil {
		return nil, err
	}

	var offset uint32
	var events []*Event
	for offset <= uint32(n-syscall.SizeofInotifyEvent) {
		ie := (*syscall.InotifyEvent)(unsafe.Pointer(&buf[offset]))
		event := &Event{
			wd:     uint32(ie.Wd),
			mask:   ie.Mask,
			cookie: ie.Cookie,
		}
		event.Path = w.wdToPathMap[event.wd]

		offset += syscall.SizeofInotifyEvent
		if ie.Len > 0 {
			nameBytes := (*[syscall.PathMax]byte)(unsafe.Pointer(&buf[offset]))
			event.Name = strings.TrimRight(string(nameBytes[0:ie.Len]), "\000")
			offset += ie.Len
		}

		events = append(events, event)
	}

	return events, nil
}

func (w *Watcher) IsUnreadEvent(event *Event) bool {
	if event.wd != w.pathToWdMap[event.Path] {
		return true
	}

	return false
}

func (w *Watcher) Free() {
	for path, _ := range w.pathToWdMap {
		w.RmWatch(path)
	}
	syscall.Close(w.fd)
}
