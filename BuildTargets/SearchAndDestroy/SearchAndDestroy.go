package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"

	"rabbitsden.online/FindDuplicateFiles/Constants"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Takes regular expression and the file name, finds all lines in the file that match the regular expression,
// interprets them as the filenames and deletes files with these names.
// Basically, this is an equivalent of
// 	 egrep <regular expression> <file> | xargs rm
// but works on the platforms where neither grep nor xargs are available (Windows, I am looking at you) and
// also deals with filenames with the special characters correctly.
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <regular expression> <file>\n", os.Args[0])
		os.Exit(1)
	}
	startTS := time.Now()
	nFilesDeleted, nFilesNotDeleted, err := findAndDeleteFiles(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error %s after deleting %d files and skipping %d files\n",
			err, nFilesDeleted, nFilesNotDeleted)
		os.Exit(1)
	}
	endTS := time.Now()
	fmt.Printf("Successfully deleted files %d and skipped %d files in %d ms (%.02f files/s) \n",
		nFilesDeleted, nFilesNotDeleted,
		endTS.Sub(startTS)/time.Millisecond, float64(nFilesDeleted)/(float64(endTS.Sub(startTS))/float64(time.Second)))
}

const (
	STATE_START            = 0
	STATE_IN_THE_BLOCK     = 1
	STATE_OUT_OF_THE_BLOCK = 2
)

func findAndDeleteFiles(pattern string, fileName string) (int, int, error) {
	rFile := regexp.MustCompile(pattern)
	rSeparator := regexp.MustCompile(`^` + Constants.BLOCK_SEPARATOR + `$`)
	// Read lines one by one
	file, err := os.Open(fileName)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()
	nDeleted := 0
	nNotDeleted := 0
	scanner := bufio.NewScanner(file)
	state := STATE_START
	nToBeRemoved := 0
	nToBeKept := 0
	var filesToRemove []string
	for scanner.Scan() {
		line := scanner.Text()
		if rSeparator.MatchString(line) {
			switch state {
			case STATE_START:
				state = STATE_IN_THE_BLOCK
				filesToRemove = make([]string, 0)
				nToBeRemoved = 0
				nToBeKept = 0
			case STATE_IN_THE_BLOCK:
				// Now we need to figure out whether we are deleting all files
				// or still leaving some
				for _, fname := range filesToRemove {
					if nToBeKept == 0 {
						fmt.Printf("Not deleting %s because there will be no such files left\n", fname)
						nNotDeleted++
					} else {
						fmt.Printf("Deleting %s\n", fname)
						err := os.Remove(fname)
						if err != nil {
							return nDeleted, nNotDeleted, err
						}
						nDeleted++
					}
				}
				state = STATE_IN_THE_BLOCK
				filesToRemove = make([]string, 0)
				nToBeRemoved = 0
				nToBeKept = 0
			}
			continue
		}
		if state == STATE_IN_THE_BLOCK {
			if rFile.MatchString(line) {
				filesToRemove = append(filesToRemove, line)
				nToBeRemoved++
			} else {
				nToBeKept++
			}
		}
	}
	return nDeleted, nNotDeleted, nil
}
