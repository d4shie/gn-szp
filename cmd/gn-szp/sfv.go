package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"os"
	"sync"
)

func modeSfvOther() {
	if gIsFileSfv {
		if gErr = loadSfv(gPath + gFilename); gErr != nil {
			os.Exit(1)
		}

		// using File.Readdirnames allows us to be up to 20x faster
		// than if we'd os.ReadDir()'d, as it doesn't have to iterate
		// over directory information, and can simply get us the names.
		var dir *os.File
		var items []string
		var lenItems int
		if dir, gErr = os.Open(gPath); gErr == nil {
			items, gErr = dir.Readdirnames(-1)
			if lenItems = len(items); lenItems == 0 && (gErr != nil && gErr != io.EOF) {
				os.Exit(1)
			}
		}

		wg := new(sync.WaitGroup)
		for _, n := range gSfvFiles {
			wg.Add(1)

			go makeMissing(n, items, wg)
		}
		wg.Wait()
		return
	}

}

// loadSfv loads and parses a .sfv file as uploaded by an end-user and loads
// the file names and hashes into gSfvFiles and gSfvHashes respectively.
func loadSfv(path string) error {
	value, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// todo: the newline splitter here needs a rework; it won't work if
	//       there's mixed \r\n and \n.
	values := bytes.Split(value, []byte("\r\n"))
	valueLen := 0
	if valueLen = len(values); valueLen == 0 {
		values = bytes.Split(value, []byte("\n"))
		if valueLen = len(values); valueLen == 0 {
			return nil
		}
	}

	// loop over once before primary, so we can adjust amount to exclude
	// newlines and comments
	for _, value = range values {
		if len(value) == 0 || value[0] == ';' {
			valueLen--
		}
	}

	gSfvFiles = make([]string, valueLen)
	gSfvHashes = make([]uint32, valueLen)

	i := 0
	for _, value = range values {
		if valueLen = len(value); valueLen == 0 || value[0] == ';' {
			continue
		}

		// HACK: we reuse valueLen here to avoid assigning another heap-allocated variable,
		//       because for some reason, the go compiler decided to allocate this int32
		//       on the fucking heap.
		valueLen = bytes.IndexByte(value, ' ')   // valueLen -> indexOfSpace
		gSfvFiles[i] = string(value[0:valueLen]) // valueLen -> indexOfSpace

		// HACK: as above, we reuse valueLen here, but this time in multiple ways!
		//       yay for order of operations! -- anyway, there's multiple hacks
		//       on this line, so let's break them down, so both my future self,
		//       and you, poor soul of a reader who's wanting to work on this,
		//       can understand this later.
		//
		//       within hex.Decode, we use `value` as both our source, and our destination,
		//       this allows us to avoid any further allocations as it will overwrite the
		//       first X bytes of the destination slice, given that the destination is
		//       large enough to fit them.
		//
		//       When providing the source, we provide an offset of value which is calculated
		//       by taking the last hack we did above ("indexOfSpace" put into valueLen),
		//       and adding the offset of the "last index byte" of rune ' ' (space),
		//       and then adding 1 (go is inclusive, so we'd get the space we found if we didn't)
		//
		//       once we've provided both, we are able to discard our mental note of "indexOfSpace" for
		//       valueLen until the next loop, so we reuse it as "bytesWritten", which gives us
		//       the amount of bytes translated from hex to the slice.
		//       We use that below in the LE uint32 conversion :)
		valueLen, err = hex.Decode(value, value[valueLen+(bytes.LastIndexByte(value[valueLen:], ' ')+1):])
		if err != nil {
			// todo: handle invalid sfv entry (either via exit(1) or continue)
			//       additionally, it's possible to check for various fail states
			//       (such as the unknown value being \r,\n), as well as if we were
			//       able to decode at least 4 bytes before the failure happened.
			//       and recover from such a state.
			panic(err)
		}

		// HACK: yet again, as above, we use valueLen here as "bytesWritten"
		gSfvHashes[i] = binary.LittleEndian.Uint32(value[:valueLen])
		i++
	}

	return nil
}

// loadCustomSfvData loads our own custom binary format for SFV information
// data format: {uint8(sfvMode)}{uint32(items)}[{uint32(crc32)}{uint32(string length (incl. null)}{null-terminated-string(filename)}...](repeated)
func loadCustomSfvData() bool {
	arrIdx := 0
	off := uint32(0)
	fileLen := uint32(0)
	items := uint32(0)
	itemLen := uint32(0)

	// todo: load the sfvdata file - we're only called if it exists
	fileBytes := []byte{}
	// uint(items) + uint(first crc32) + uint(first string length)
	if len(fileBytes) < 12 {
		return false
	}

	items = binary.LittleEndian.Uint32(fileBytes)
	gSfvFiles = make([]string, items)
	gSfvHashes = make([]uint32, items)

	for off <= fileLen {
		gSfvHashes[arrIdx] = binary.LittleEndian.Uint32(fileBytes[off : off+4])
		off += 4

		itemLen = binary.LittleEndian.Uint32(fileBytes[off : off+4])
		off += 4
		if itemLen == 0 {
			continue
		}

		gSfvFiles[arrIdx] = string(fileBytes[off : off+itemLen])
		off += itemLen
	}

	return true
}
