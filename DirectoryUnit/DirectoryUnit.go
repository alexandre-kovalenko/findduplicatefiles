package DirectoryUnit

import (
	"errors"
	"fmt"
	"log"
	"os"

	"rabbitsden.online/FindDuplicateFiles/FileUnit"
)

type DirectoryUnit struct {
	Name       string
	PlainFiles []FileUnit.FileUnit
}

const IGNORE_DOT_FILES = false

// /////////////////////////////////////////////////////////////////////////////
// Creates collection of directory units for the specified directory and all
// subdirectories.
func MakeDirectoryUnits(name string) (map[string]DirectoryUnit, error) {
	result := make(map[string]DirectoryUnit)
	myUnit := DirectoryUnit{Name: name}
	//////////////////////////////////////////////////////////////////////////
	//
	entries, err := os.ReadDir(name)
	if err != nil {
		message := fmt.Sprintf("Error reading directory %s (%v)", name, err)
		log.Println(message)
		return nil, errors.New(message)
	}
	for _, e := range entries {
		if IGNORE_DOT_FILES && e.Name()[0] == '.' {
			continue
		}
		if e.IsDir() {
			subName := name + "/" + e.Name()
			log.Printf("Recursing into %s\n", subName)
			subUnits, err := MakeDirectoryUnits(subName)
			if err != nil {
				return nil, err
			}
			for k, v := range subUnits {
				result[name+"/"+k] = v
			}
		} else {
			fullPath := name + "/" + e.Name()
			fu, err := FileUnit.MakeFileUnit(fullPath)
			if err != nil {
				message := fmt.Sprintf("Error reading file %s (%v)", fullPath, err)
				log.Println(message)
				return nil, errors.New(message)
			}
			myUnit.PlainFiles = append(myUnit.PlainFiles, fu)
		}
	}
	//
	///////////////////////////////////////////////////////////////////////////
	result[name] = myUnit
	return result, nil
}
