package inotify

import "syscall"

const (
	IN_ALL_EVENTS = syscall.IN_ALL_EVENTS
	IN_IGNORED    = syscall.IN_IGNORED

	IN_MODIFY = syscall.IN_MODIFY
	IN_ATTRIB = syscall.IN_ATTRIB

	IN_MOVE_SELF  = syscall.IN_MOVE_SELF
	IN_MOVED_FROM = syscall.IN_MOVED_FROM
	IN_MOVED_TO   = syscall.IN_MOVED_TO

	IN_DELETE_SELF = syscall.IN_DELETE_SELF
	IN_DELETE      = syscall.IN_DELETE

	IN_CREATE = syscall.IN_CREATE
)

type Event struct {
	wd     uint32
	mask   uint32
	cookie uint32

	Path string
	Name string
}

func (e *Event) InIgnored() bool {
	return e.mask&IN_IGNORED == IN_IGNORED
}

func (e *Event) InModify() bool {
	return e.mask&IN_MODIFY == IN_MODIFY
}

func (e *Event) InAttrib() bool {
	return e.mask&IN_ATTRIB == IN_ATTRIB
}

func (e *Event) InMoveSelf() bool {
	return e.mask&IN_MOVE_SELF == IN_MOVE_SELF
}

func (e *Event) InMovedFrom() bool {
	return e.mask&IN_MOVED_FROM == IN_MOVED_FROM
}

func (e *Event) InMovedTo() bool {
	return e.mask&IN_MOVED_TO == IN_MOVED_TO
}

func (e *Event) InDeleteSelf() bool {
	return e.mask&IN_DELETE_SELF == IN_DELETE_SELF
}

func (e *Event) InDelete() bool {
	return e.mask&IN_DELETE == IN_DELETE
}

func (e *Event) InCreate() bool {
	return e.mask&IN_CREATE == IN_CREATE
}
