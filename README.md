This is the set of utilities to remove duplicate files from single or multiple trees. Duplicate files are defined 
as having same SHA-256 hash.

There are three utilities:
1. FindDuplicateFiles -- takes an optional output file name and the list of directories to look for duplicates.
   If output file name is not specified, STDOUT is used.
2. SearchAndDestroy -- takes regular expression and the file, created by FindDuplicateFiles and removes all duplicates
   that match regular expression. If all files in the duplicate group match the regular expression, no files are removed.
3. PruneEmptyDirectories -- takes single path and removes all directories underneath that do not contain plain files.

All three utilities could be built by changing directory to BuildTargets/<Utility name> and issuing 'go build'.
