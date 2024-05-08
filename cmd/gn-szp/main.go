package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"syscall"
)

func main() {
	// The following is an explanation of the required arguments, environment
	// variables, and other requirements for initialization:
	//
	// Args:
	//   gn-szp <filename> <path> <crc>
	//     filename: the name of the file that was uploaded for this invocation
	//     path:     the path the file was uploaded to
	//     crc:      hex-encoded on-the-fly crc calculation by glftpd;
	//                   this is controlled by the "calc_crc" config setting
	//                   and if the option is disabled or not applicable,
	//                   we will receive the value "00000000".
	//
	// Environment variables:
	//   $USER:     the username of the authenticated user
	//   $GROUP:    the primary group of the authenticated user
	//                  if the group is empty, defaults to "NoGroup"
	//   $TAGLINE:  the tagline of the authenticated user
	//   $SPEED:    the average upload/download speed of the authenticated user
	//   $SECTION:  the section the authenticated user is currently in
	//                  if the section is empty, defaults to "DEFAULT"
	if len(os.Args) < 4 {
		fmt.Printf("gn-szp: wrong number of arguments used\n"+
			"Usage: %s <filename> <path> <crc>\n", os.Args[0])
		os.Exit(1)
	}

	gFilename = os.Args[1]
	gPath = os.Args[2]
	gOTFCrc = os.Args[3]

	// here we use syscall.Getenv instead of os.Getenv because we want to skip the
	// extra calls os.Getenv makes (see pkg testlog)
	gUser, _ = syscall.Getenv("USER")
	if gGroup, gOk = syscall.Getenv("GROUP"); !gOk {
		gGroup = "NoGroup"
	}
	if gSection, gOk = syscall.Getenv("SECTION"); !gOk {
		gSection = "DEFAULT"
	}
	gSpeed, _ = syscall.Getenv("SPEED")
	gTagline, _ = syscall.Getenv("TAGLINE")

	if _, gErr = os.Stat(DataPath + gPath); os.IsNotExist(gErr) {
		// if this is a completely new entry that we have no metadat for
		// create it, and proceed to chooseRunner directly

		gErr = os.MkdirAll(DataPath+gPath, 0777)
		if gErr != nil {
			panic(gErr)
		}
	} else if _, gErr = os.Stat(DataPath + gPath + "sfvdata"); gErr == nil {
		// if we have metadata for it, and we've previously parsed the SFV, load
		// the sfv data file, so we can have a clue as to what's going on

		// todo: load sfvdata
		gSfvDataLoaded = loadCustomSfvData()
		// todo: handle cases where sfvdata file was not loaded
		//       e.g: malformed sfvdata, if .sfv exists, scan it first, then retry
		//            if no sfv exists, and the SfvFirst option is enabled, reject file
		// todo: have other supported modes like SFV+VIDEO, SFV+AUDIO, SFV+RAR
		gRunMode = RunModeSfvOther
	} // TODO: Support ZIP+DIZ

	chooseRunner()
}

func makeMissing(fname string, itemsAtTimeOfCheck []string, wg *sync.WaitGroup) {
	gConcurrentLimiter <- struct{}{}
	defer func() {
		<-gConcurrentLimiter
		wg.Done()
	}()

	var f *os.File
	// TODO: Look into adding optional case-insensitivity for SFV files, which
	//       if enabled uses EqualFold in slices.ContainsFunc
	//       WARNING: strings.Equals is NOT faster than using the "==" operator,
	//                but EqualFold IS FASTER THAN THE strings.ToLower() method,
	//                as it checks rune-by-rune for codepoints, cutting off short
	//                on the first mismatch.
	if !slices.Contains(itemsAtTimeOfCheck, fname) {
		// Slow path using filesystem-backed check to ensure that we won't
		// accidentally make a -missing file during a race.
		if _, err := os.Stat(gPath + fname); os.IsNotExist(err) {
			f, err = os.Create(gPath + fname + MissingSuffix)
			if err != nil {
				gLog.Debug("failed to create -missing file: %v", err)
				return
			}
			defer f.Close()
		}
	}
}

// chooseRunner determines the function to use for the current gRunMode
func chooseRunner() {
	switch gRunMode {
	case RunModeSfvOther:
		modeSfvOther()
		break
	case RunModeUndetermined:
		fallthrough
	default:
		determineRunMode()
		break
	}
}

// determineRunMode is called when metadata does not exist for a given release
// it uses the filename, filepath, and directory structure information to try
// and determine what type of release this is, and set gRunMode accordingly
//
// if gRunMode is set, chooseRunner is called to start the process
func determineRunMode() {
	// only implementing SFV+OTHER for now
	for _, dir := range gSfvDirs {
		if strings.HasPrefix(gPath, dir) && strings.HasSuffix(gFilename, ".sfv") {
			gRunMode = RunModeSfvOther
			gIsFileSfv = true
			break
		}
	}

	if gRunMode == RunModeUndetermined {
		//panic("run mode could not be determined, this should be turned into a fatal loggable instead of a panic")
		os.Exit(1)
	}

	chooseRunner()
}
