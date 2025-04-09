package main

import (
	"log"
	"os"

	"rabbitsden.online/FindDuplicateFiles/DirectoryUnit"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("Usage: %s <directory>\n", os.Args[0])
		return
	}
	rootDirectory := os.Args[1]
	dUnits, err := DirectoryUnit.MakeDirectoryUnits(rootDirectory)
	if err != nil {
		log.Printf("Failed to enumerate %s (%v)\n", rootDirectory, err)
		os.Exit(-1)
	}
	csMap := make(map[string][]string)
	for _, du := range dUnits {
		for _, f := range du.PlainFiles {
			eSum := f.GetEncodedChecksum()
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
}

// //////////////////////////////////////////////////////////////////////////////////
// Process duplicate files
func processDuplicates(duplicates []string) error {
	for _, f := range duplicates {
		log.Printf("%s\n", f)
	}
	return nil
}
