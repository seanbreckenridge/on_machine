package on_machine

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func matchIgnoreExtensions(currentTarget string, base string, pathName string, ext string) []string {
	matches := make([]string, 0)
	// since the currentTarget (pattern) cant include an extension, since we extracted it into
	// ext in MatchPaths, this is matching a directory
	baseWithExt := filepath.Join(base, pathName)
	if pathName == currentTarget {
		if exists, _ := PathExists(baseWithExt); exists {
			matches = append(matches, baseWithExt)
		}
	}
	// try to remove the extension from the file
	fileExt := filepath.Ext(pathName)
	// see with_extensions example structure for why this is needed
	// i.e., to include both the linux directory, and the linux.zsh file
	pathWithoutExt := strings.TrimSuffix(pathName, fileExt)
	if pathWithoutExt == currentTarget {
		baseWithoutExt := filepath.Join(base, pathWithoutExt)
		if exists, _ := PathExists(baseWithoutExt); exists {
			matches = append(matches, baseWithoutExt)
		}
		if ext != "" {
			// try adding the extracted extension from the end of the pattern here
			baseAddedExt := filepath.Join(base, pathWithoutExt+ext)
			if exists, _ := PathExists(baseAddedExt); exists {
				matches = append(matches, baseAddedExt)
			}
		}
	}
	return SliceUniqMap(matches)
}

func matchRecurHelper(targetsRendered []string, currentDepth int, currentBase string, ext string) ([]string, error) {
	matches := make([]string, 0)
	// base cases
	// the current base (joined from above recursive call doesn't exist, so we can't search further recursively)
	if !DirExists(currentBase) {
		return matches, nil
	}
	if currentDepth > len(targetsRendered)-1 {
		// if we cant index into the targetsRendered anymore, that means we've passed the depth this allows
		return matches, nil
	}
	currentTarget := targetsRendered[currentDepth]
	dirContents, err := os.ReadDir(currentBase)
	if err != nil {
		return matches, err
	}
	for _, pthInfo := range dirContents {
		pthName := pthInfo.Name()
		// all acts as a wildcard ('*'), it means items in that directory should always be included
		// should match with or without extension, to allow for nested structures
		for _, litPattern := range []string{currentTarget, "all"} {
			// try to match this pattern with/without the extension
			for _, match := range matchIgnoreExtensions(litPattern, currentBase, pthName, ext) {
				matches = append(matches, match)
				// try to match items recursively at a depth lower
				// base case returns if fullPth isn't a directory/doesn't exist
				recurMatches, err := matchRecurHelper(targetsRendered, currentDepth+1, match, ext)
				if err != nil {
					return matches, err
				}
				matches = append(matches, recurMatches...)
			}
		}
	}
	return matches, nil
}

func MatchPaths(pattern string, baseDir string) ([]string, error) {
	matches := make([]string, 0)
	if strings.TrimSpace(pattern) == "" {
		return matches, errors.New("matching pattern can not be an empty string")
	}
	// if there is an extension on the pattern, separate it
	ext := filepath.Ext(pattern)
	if ext != "" {
		pattern = strings.TrimSuffix(pattern, ext)
	}

	//
	// spit the pattern into directory parts, if possible
	// don't use OS-specific path separator, whole point of this
	// is it works cross platform, so leaving it as '/' means
	// it doesn't have be changed up bash/above cmdline argument
	targetsRaw := strings.Split(pattern, "/")
	targetsRendered := make([]string, len(targetsRaw))
	// for each directory part, render to expected value
	for i, part := range targetsRaw {
		rendered := ReplaceFields(part)
		targetsRendered[i] = rendered
	}
	// convert to abspath
	base, err := filepath.Abs(baseDir)
	if err != nil {
		return matches, err
	}
	newMatches, err := matchRecurHelper(targetsRendered, 0, base, ext)
	if err != nil {
		return matches, err
	}
	matches = append(matches, newMatches...)
	matches = SliceUniqMap(matches)
	return matches, nil
}
