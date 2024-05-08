// By default, we disable Go's garbage collector as gn-szp is not meant to
// be run as a long-running process, but rather as a hook on every file
// uploaded, so we don't have a need to GC most of the time.
//
// It's possible to enable the GC in the case of memory-usage issues
// by compiling with the enable_gc build tag.
//go:build !enable_gc

package main

import "runtime/debug"

func init() {
	debug.SetGCPercent(-1)
}
