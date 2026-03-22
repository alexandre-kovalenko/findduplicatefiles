package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"rabbitsden.online/FindDuplicateFiles/Constants"
	"rabbitsden.online/FindDuplicateFiles/DirectoryUnit"
	"rabbitsden.online/FindDuplicateFiles/UnitLimiter"
)

var resultFilename string

func Usage() {
	log.Printf("Usage: %s [-o <result filename>] <directory> [<directory>...]\n", os.Args[0])
	os.Exit(-1)
}

func main() {
	if len(os.Args) < 2 {
		Usage()
	}
	if (os.Args[1] == "-o") && (len(os.Args) < 4) {
		Usage()
	}
	firstDirectoryIndex := 1
	if os.Args[1] == "-o" {
		resultFilename = os.Args[2]
		_ = os.Remove(resultFilename)
		firstDirectoryIndex = 3
	}
	// Since we are heavily I/O bound, let's schedule 8 goroutines per core
	runtime.GOMAXPROCS(runtime.NumCPU() * 8)
	log.Printf("Running on %d cores, max procs is: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(0))
	startTS := time.Now()
	fl := UnitLimiter.MakeUnitLimiter(1000, "file")
	csMap := make(map[string][]string)
	filesProcessed := 0
	for i := firstDirectoryIndex; i < len(os.Args); i++ {
		rootDirectory := os.Args[i]
		mdu := DirectoryUnit.MakeDirectoryUnits(rootDirectory, &fl)
		dUnits, err := mdu.DirectoryUnits, mdu.Error
		if err != nil {
			log.Printf("Failed to enumerate %s (%v)\n", rootDirectory, err)
			os.Exit(-1)
		}
		for _, du := range dUnits {
			for _, f := range du.PlainFiles {
				eSum := f.GetEncodedChecksum()
				filesProcessed++
				if _, ok := csMap[eSum]; ok {
					csMap[eSum] = append(csMap[eSum], f.Name)
				} else {
					csMap[eSum] = []string{f.Name}
				}
			}
		}
	}
	for _, v := range csMap {
		if len(v) > 1 {
			err := processDuplicates(v)
			if err != nil {
				log.Printf("Failed to process duplicates in %v -- (%v)\n", v, err)
				os.Exit(-1)
			}
			log.Println("==================")
		}
	}
	endTS := time.Now()
	elapsedMS := endTS.Sub(startTS).Milliseconds()
	log.Printf("Running on %d cores, max procs is: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(0))
	log.Printf("Processed %d files in %d ms (%.2f files/s)\n", filesProcessed, elapsedMS,
		float64(filesProcessed)/(float64(elapsedMS)/1000))
}

// //////////////////////////////////////////////////////////////////////////////////
// Process duplicate files

func processDuplicates(duplicates []string) error {
	separateResultFile := false
	var fh *os.File
	var err error
	if len(resultFilename) > 0 {
		fh, err = os.OpenFile(resultFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		separateResultFile = true
		fmt.Fprintln(fh, Constants.BLOCK_SEPARATOR)
	}
	defer func() {
		if separateResultFile {
			fh.Close()
		}
	}()
	for _, f := range duplicates {
		log.Println(f)
		if separateResultFile {
			fmt.Fprintf(fh, "%s\n", f)
		}
	}
	return nil
}
