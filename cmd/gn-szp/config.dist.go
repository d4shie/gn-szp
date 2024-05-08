// This file is meant to serve as an "official" version of the gn-szp config.
// Please copy this file to config.go, and remove the following line before
// editing.
//go:build NEVER

package main

const (
	// SiteName is the name of your site. Used in the [sn] cookie.
	SiteName = ""
	// ShortSiteName is the 2-4 letter name for your site. Used in the [sns] cookie.
	ShortSiteName = ""

	// SitePath is where your site's data is housed under.
	SitePath = "/site/"
	// DataPath is the location where metadata is stored.
	// It serves the exact same purpose as pzs' "storage" setting.
	DataPath = "/ftp-data/gn-szp/"
	// LocksPath is the location where lock files are stored.
	// lock files are used to ensure that the metadata in a race can only be
	// modified by a single instance of the program.
	//
	// It is recommended that this path is mounted as a tmpfs/shmfs/ramdisk.
	LocksPath = "/tmp/gn-szp/"

	// SfvDirs is a space-separated list of directories that will use SFV modes
	SfvDirs = ""

	GroupDirs = ""

	// IgnoreDirs is a space-separated list of directories that will be ignored.
	// Same as nocheck_dirs in pzs.
	IgnoreDirs = ""

	// MaximumConcurrency is the max amount of threads (goroutines, in this case)
	// that may be launched at one time.
	//
	// Please note that this is heavily hardware-dependent, and requires that
	// you size it up correctly, otherwise you may run into performance issues
	// with multiple files being processed at the same time.
	MaximumConcurrency = 300

	// MissingSuffix is the suffix of the meta-files that are created in SFV mode
	// if an item has not been uploaded yet.
	MissingSuffix = "-missing"
)

// WARNING: Do not edit!
// By generating the slices ahead of time, we are able to shave off a lot of
// startup time (in the tens of ms!) when there are many directories defined.
//
// As of at least Go 1.22, when there is a static-length slice that is not
// modified throughout the lifetime of the program, it is emitted as a constant
// value in the data section, so we don't even incur an initialization penalty!

//go:generate go run ../dirsplit_gen ./.dummy -clean
//go:generate go run ../dirsplit_gen ./.SfvDirs gSfvDirs
//go:generate go run ../dirsplit_gen ./.GroupDirs gGroupDirs
//go:generate go run ../dirsplit_gen ./.IgnoreDirs gIgnoreDirs
