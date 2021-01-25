package main

//File Fiesta takes a directory (dir) and number (numFiles) input from the user
//It runs a recursive depth first search (using filepath.Walk) through all underlying files and folders in the Dir
//Return a sorted list of the largest N files, with file name, location, & size

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/text/message"
)

type file struct {
	name     string
	location string
	size     int64
}

func main() {
	//get directory & number of files to return from command line flags
	dir := flag.String("dir", "./", "Directory to search")
	numFiles := flag.Int("numFiles", 20, "Number of files to return")
	flag.Parse()

	//create a new printer that formats numbers using 000s commas
	p := message.NewPrinter(message.MatchLanguage("en"))

	//File Fiesta!!
	p.Println("\n---------------------------------File Fiesta---------------------------------")

	//call the fileSearch function and save the result
	results, fileCount, skippedDirs, dirSize, err := fileSearch(*dir, *numFiles)
	if err != nil {
		panic(err)
	}

	//print the results to the terminal
	p.Println("The subject directory is", dirSize/1000000, "mb, including any hidden directories.")
	p.Println("\n\nTotal files searched:\t", fileCount)
	p.Println("Hidden directories skipped:\t", skippedDirs)
	if fileCount < *numFiles {
		p.Println("\nReturned the largest", fileCount, "files:")
	} else {
		p.Println("\nReturned the largest", *numFiles, "files:")
	}
	for i, s := range results {
		if s.name == "." {
			continue
		}
		p.Println(i+1, "------------------------------------------------")
		p.Println("Name:\t\t", s.name)
		p.Println("Location:\t", s.location)
		p.Println("Size\t\t", s.size/1000000, "mb\n")
	}

}

func fileSearch(dir string, numFiles int) ([]file, int, int, int64, error) {
	var dirSize int64
	var smallest int64
	// var smallIndex int
	fileCount := 0
	skippedDirs := 0
	files := []file{}

	//filepath.Walk does a depth first, recursive search through a directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//skip hidden directories (if chosen in the flags)
		if info.IsDir() && info.Name() != filepath.Base(dir) && info.Name()[0] == '.' {
			skippedDirs++
			dirSize += info.Size()
			return filepath.SkipDir
		}

		//only include the largest files (numFiles)
		if !info.IsDir() {
			if info.Size() > smallest && len(files) < numFiles {
				files = append(files, file{name: info.Name(), location: path, size: info.Size()})
			}
		}

		//count all all locations except the root path we are searching
		if info.Name() != filepath.Base(dir) {
			fileCount++
		}

		//cummulative size of all files we've visited
		dirSize += info.Size()

		return nil
	})
	if err != nil {
		fmt.Println("There was an error while walking the directory:", err)
		return nil, fileCount, skippedDirs, dirSize, err
	}

	return files, fileCount, skippedDirs, dirSize, nil
}
