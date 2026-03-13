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
func MakeDirectoryUnits(name string, ch chan MakeDirectoryUnitsResult) {
	dirCallList := make([]DirCallList, 0)
	fileCallList := make([]FileCallList, 0)
	result := make(map[string]DirectoryUnit)
	myUnit := DirectoryUnit{Name: name}
	//////////////////////////////////////////////////////////////////////////
	//
	entries, err := os.ReadDir(name)
	if err != nil {
		message := fmt.Sprintf("Error reading directory %s (%v)", name, err)
		log.Println(message)
		ch <- MakeDirectoryUnitsResult{nil, errors.New(message)}
	}
	for _, e := range entries {
		if IGNORE_DOT_FILES && e.Name()[0] == '.' {
			continue
		}
		if e.IsDir() {
			subCh := make(chan MakeDirectoryUnitsResult)
			subName := name + "/" + e.Name()
			log.Printf("Recursing into %s\n", subName)
			dirCallList = append(dirCallList, DirCallList{name, subCh})
			go MakeDirectoryUnits(subName, subCh)
		} else {
			fullPath := name + "/" + e.Name()
			subCh := make(chan FileUnit.MakeFileUnitResult)
			fileCallList = append(fileCallList, FileCallList{fullPath, subCh})
			go FileUnit.MakeFileUnit(fullPath, subCh)
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
			ch <- MakeDirectoryUnitsResult{nil, errors.New(message)}
		}
		myUnit.PlainFiles = append(myUnit.PlainFiles, fu)
		close(cl.ch)
	}
	///////////////////////////////////////////////////////////////////////////
	// At this point goroutines for every subdirectory have been started, so
	// we need to wait for them and fold the results into overall result.
	for _, cl := range dirCallList {
		subMDU := <-cl.ch
		if subMDU.Error != nil {
			ch <- MakeDirectoryUnitsResult{nil, subMDU.Error}
		}
		for k, v := range subMDU.DirectoryUnits {
			result[k] = v
		}
		close(cl.ch)
	}
	//
	///////////////////////////////////////////////////////////////////////////
	result[name] = myUnit
	ch <- MakeDirectoryUnitsResult{result, nil}
}
