package main

// RunMode is the operating mode used by gn-szp for a given execution.
// It determines the execution path and what checks will be run, as well as the
// metadata that will be generated.
type RunMode byte

const (
	RunModeUndetermined RunMode = iota
	RunModeZipDiz
	RunModeSfvAudio
	RunModeSfvVideo
	RunModeSfvRar
	RunModeSfvOther
)
