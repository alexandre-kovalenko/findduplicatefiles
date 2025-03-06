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
	for _, du := range dUnits {
		for _, f := range du.PlainFiles {
			log.Printf("%s\n", f.ToString())
		}
	}
}
