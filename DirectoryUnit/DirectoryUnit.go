package DirectoryUnit

import (
	"errors"
	"fmt"
	"log"
	"os"

	"rabbitsden.online/FindDuplicateFiles/FileUnit"
	"rabbitsden.online/FindDuplicateFiles/UnitLimiter"
)

type DirectoryUnit struct {
	Name       string
	PlainFiles []FileUnit.FileUnit
}

type MakeDirectoryUnitsResult struct {
	DirectoryUnits map[string]DirectoryUnit
	Error          error
}

type DirCallList struct {
	Name string
	ch   chan MakeDirectoryUnitsResult
}

type FileCallList struct {
	Name string
	ch   chan FileUnit.MakeFileUnitResult
}

const IGNORE_DOT_FILES = false

// /////////////////////////////////////////////////////////////////////////////
// Creates collection of directory units for the specified directory and all
// subdirectories.
func MakeDirectoryUnits(name string, fl *UnitLimiter.Limiter) MakeDirectoryUnitsResult {
	fileCallList := make([]FileCallList, 0)
	result := make(map[string]DirectoryUnit)
	myUnit := DirectoryUnit{Name: name}
	//////////////////////////////////////////////////////////////////////////
	//
	entries, err := os.ReadDir(name)
	if err != nil {
		message := fmt.Sprintf("Error reading directory %s (%v)", name, err)
		log.Println(message)
		return MakeDirectoryUnitsResult{nil, errors.New(message)}
	}
	for _, e := range entries {
		if IGNORE_DOT_FILES && e.Name()[0] == '.' {
			continue
		}
		if e.IsDir() {
			subName := name + "/" + e.Name()
			log.Printf("Recursing into %s\n", subName)
			subMDU := MakeDirectoryUnits(subName, fl)
			if subMDU.Error != nil {
				return MakeDirectoryUnitsResult{nil, subMDU.Error}
			}
			for k, v := range subMDU.DirectoryUnits {
				result[k] = v
			}
		} else {
			fullPath := name + "/" + e.Name()
			subCh := make(chan FileUnit.MakeFileUnitResult)
			fileCallList = append(fileCallList, FileCallList{fullPath, subCh})
			go FileUnit.MakeFileUnit(fullPath, subCh, fl)
		}
	}
	///////////////////////////////////////////////////////////////////////////
	// First we collect data from all of the plain files
	for _, cl := range fileCallList {
		fullPath := cl.Name
		r := <-cl.ch
		fu, err := r.FileUnit, r.Error
		if err != nil {
			message := fmt.Sprintf("Error reading file %s (%v)", fullPath, err)
			log.Println(message)
			return MakeDirectoryUnitsResult{nil, errors.New(message)}
		}
		myUnit.PlainFiles = append(myUnit.PlainFiles, fu)
		close(cl.ch)
	}
	result[name] = myUnit
	return MakeDirectoryUnitsResult{result, nil}
}
