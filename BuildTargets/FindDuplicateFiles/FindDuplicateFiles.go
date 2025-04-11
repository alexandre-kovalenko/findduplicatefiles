package main

import (
	"log"
	"os"
	"time"

	"rabbitsden.online/FindDuplicateFiles/DirectoryUnit"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("Usage: %s <directory>\n", os.Args[0])
		return
	}
	startTS := time.Now()
	rootDirectory := os.Args[1]
	ch := make(chan DirectoryUnit.MakeDirectoryUnitsResult)
	go DirectoryUnit.MakeDirectoryUnits(rootDirectory, ch)
	mdu := <-ch
	dUnits, err := mdu.DirectoryUnits, mdu.Error
	if err != nil {
		log.Printf("Failed to enumerate %s (%v)\n", rootDirectory, err)
		os.Exit(-1)
	}
	csMap := make(map[string][]string)
	filesProcessed := 0
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
	log.Printf("Processed %d files in %d ms (%.2f files/s)\n", filesProcessed, elapsedMS,
		float64(filesProcessed)/(float64(elapsedMS)/1000))
}

// //////////////////////////////////////////////////////////////////////////////////
// Process duplicate files
const PREFIX = "/eBooks/eBooks"

func processDuplicates(duplicates []string) error {
	for _, f := range duplicates {
		// Remove the duplicates with the given prefix
		// if len(f) >= len(PREFIX) && f[0:len(PREFIX)] == PREFIX {
		// 	log.Printf("Removing %s\n", f)
		//  err := os.Remove(f)
		//  if err != nil {
		// 	 	return err
		// 	}
		// }
		log.Println(f)
	}
	return nil
}
