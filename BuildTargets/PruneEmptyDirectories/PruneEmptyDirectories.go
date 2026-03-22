package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

////////////////////////////////////////////////////////////////////////////////////////////
// Find and remove all empty directories starting from the given one
////////////////////////////////////////////////////////////////////////////////////////////

func main() {
	if len(os.Args) != 2 {
		log.Printf("Usage: %s <path>\n", os.Args[0])
		os.Exit(1)
	}
	startTS := time.Now()
	nPruned, _, err := PruneEmptyDirectories(os.Args[1])
	if err != nil {
		log.Printf("Failed to prune empty directories (%v)... %d have been puned prior to error.\n",
			err, nPruned)
		os.Exit(1)
	}
	endTS := time.Now()
	log.Printf("Successfully pruned %d empty directories in %d ms (%.02f directories/s) \n",
		nPruned,
		endTS.Sub(startTS)/time.Millisecond, float64(nPruned)/(float64(endTS.Sub(startTS))/float64(time.Second)))
}

const IGNORE_DOT_FILES = false

func PruneEmptyDirectories(startingPath string) (int, int, error) {
	pruned := 0
	survivingChildren := 0
	entries, err := os.ReadDir(startingPath)
	if err != nil {
		message := fmt.Sprintf("Error reading directory %s (%v)", startingPath, err)
		log.Println(message)
		return pruned, 1, errors.New(message)
	}
	for _, e := range entries {
		if IGNORE_DOT_FILES && e.Name()[0] == '.' {
			continue
		}
		if e.IsDir() {
			subName := startingPath + "/" + e.Name()
			log.Printf("Recursing into %s\n", subName)
			nPruned, nSurvivedChildren, err := PruneEmptyDirectories(subName)
			pruned += nPruned
			if err != nil {
				message := fmt.Sprintf("Error pruning empty directories under %s (%v)", subName, err)
				return pruned, 1, errors.New(message)
			}
			if nSurvivedChildren == 0 {
				err = os.Remove(subName)
				if err != nil {
					log.Printf("Error removing empty directory %s (%v)\n", subName, err)
					survivingChildren++
				} else {
					pruned++
				}
			} else {
				survivingChildren++
			}
		} else {
			survivingChildren++
		}
	}
	return pruned, survivingChildren, nil
}
