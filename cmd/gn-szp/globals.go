package main

import "gn-szp/logging"

// In addition to the globals defined here, there are several others sourced from
// the configuration during compile-time source generation.
// Please see config.go, config.g.stub.go, and config.g.go for more information.
var (
	gRunMode RunMode         = RunModeUndetermined
	gLog     logging.ILogger = logging.Setup()

	// gSfvFiles when operating in any SFV mode, the sfvLoad function will
	// populate this global slice with the file names in the loaded SFV, in
	// the order that they appear.
	// The indices for gSfvFiles and gSfvHashes will always match.
	gSfvFiles []string
	// gSfvHashes when operating in any SFV mode, the sfvLoad function will
	// populate this global slice with the checksums in the loaded SFV, in the
	// order that they appear.
	// The indices for gSfvFiles and gSfvHashes will always match.
	gSfvHashes []uint32

	// gFilename is the filename as provided by glftpd in args
	gFilename string
	// gPath is the dir the file was uploaded to as provided by glftpd in args
	gPath string
	// gOTFCrc is the on-the-fly-calculated crc as provided by glftpd in args
	gOTFCrc string

	gUser    string // gUser is the currently-authenticated user's username.
	gGroup   string // gGroup is the currently-authenticated user's primary group.
	gTagline string // gTagline is the currently-authenticated user's tagline.
	gSpeed   string // gSpeed is the currently-authenticated user's average up/download speed.
	gSection string // gSection is the section the currently-authenticated user is in.

	// gOk is a global OK value used by synchronous-access code
	gOk bool
	// gErr is a global error value used by synchronous-access code
	gErr error

	// gIsFileSfv is set to indicate whether the accepted file was a .sfv
	gIsFileSfv bool

	// gSfvDataLoaded is set to indicate whether the loadCustomSfvData function worked or not.
	gSfvDataLoaded bool

	// gConcurrentLimiter is a global concurrency limiter that allows
	// us to prevent thread exhaustion on IO-intensive operations.
	//
	// Implementation details:
	// The concurrent limiter channel takes an empty struct due to its memory
	// footprint being literally zero. The only memory overhead we have is
	// the current length of the channel and its capacity.
	//
	// When the channel is full, the function will idly spin while waiting for
	// space in the channel to be freed so its token can be placed.
	//
	// Usage:
	// On function entry, send a value into the channel:
	//    `gConcurrentLimiter <- struct{}{}`
	//
	// On function exit (preferably, deferred), receive a value from the channel:
	//    `<-gConcurrentLimiter`
	gConcurrentLimiter = make(chan struct{}, MaximumConcurrency)
)
