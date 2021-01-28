package main

/*
File Fiesta returns a sorted list of largest files in a given folder.
It takes a string (dir) and int (numFiles) and runs a recursive,
depth first search (using filepath.Walk). It then returns the
name, path, and size of the largest files.
*/

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/text/message"
)

type file struct {
	name     string
	location string
	size     int64
}

func main() {
	timer := time.Now()
	var dir string
	var numFiles int
	var scanHidden bool

	//get directory, number of files to return, and whether to skip hidden directories from command line flags
	flag.StringVar(&dir, "dir", "./", "Directory to search?")
	flag.IntVar(&numFiles, "numFiles", 20, "Number of files to return?")
	flag.BoolVar(&scanHidden, "hidden", false, "Scan hidden directories? (default false")
	flag.Parse()
	if numFiles <= 0 {
		fmt.Println("\n**The numFiles input must be a positive integer. Please try again.**")
		os.Exit(1)
	}

	//create a new printer that formats numbers using 000s commas
	p := message.NewPrinter(message.MatchLanguage("en"))

	//File Fiesta!!
	p.Println("\n---------------------------------File Fiesta---------------------------------")

	//call the fileSearch function and save the result
	results, fileCount, skippedDirs, dirSize, err := fileSearch(scanHidden, dir, numFiles)
	if err != nil {
		p.Println("There was an error while walking the directory:")
		p.Println(err)
		p.Println("File paths that includes spaces should be surrounded by double quotes(\"\")")
		os.Exit(1)
	}

	//print the results to the terminal
	p.Printf("The subject directory is %.2f MB, including any hidden directories.\n", float64(dirSize)/1048576)
	p.Println("\n\nTotal files searched:\t\t", fileCount)
	p.Println("Hidden directories skipped:\t", skippedDirs)
	p.Println("\nReturned the largest", min([]int{fileCount, numFiles, len(results)}), "files:")

	var topSize int64
	for i, s := range results {
		if s.name == "." {
			continue
		}
		p.Println(i+1, "------------------------------------------------")
		p.Println("Name:\t\t", s.name)
		p.Println("Location:\t", s.location)
		p.Printf("Size\t\t %.2f MB\n", float64(s.size)/1048576)
		topSize += s.size
	}
	p.Println("-------------------------------Search Completed------------------------------")
	p.Println("Your search was completed in", time.Since(timer))
	p.Printf("\nThe top %d results total %.2f MB\n", min([]int{fileCount, numFiles, len(results)}), float64(topSize)/1048576)
}

func fileSearch(scanHidden bool, dir string, numFiles int) ([]file, int, int, int64, error) {
	var dirSize int64
	var smallest int64

	fileCount := 0
	skippedDirs := 0
	files := []file{}

	//filepath.Walk does a depth first, recursive search through a directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//skip hidden directories (if chosen in the flags)
		if scanHidden == false && info.IsDir() && info.Name() != filepath.Base(dir) && info.Name()[0] == '.' {
			skippedDirs++
			dirSize += info.Size()
			return filepath.SkipDir
		}

		//only include the largest files (numFiles)
		if !info.IsDir() {
			newFile := file{name: info.Name(), location: path, size: info.Size()}
			files = sortSearch(newFile, files, &smallest, numFiles)

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
		return nil, fileCount, skippedDirs, dirSize, err
	}

	return files, fileCount, skippedDirs, dirSize, nil
}

func sortSearch(newFile file, files []file, smallest *int64, numFiles int) []file {

	//if it's the first file, add it to our files slice
	if len(files) == 0 {
		files = append(files, newFile)
		*smallest = newFile.size
	} else {
		//otherwise binary search through files slice
		for i := 0; i < len(files); i++ {
			//if newFile is smaller the files[i] continue loop
			if files[i].size >= newFile.size {

				if i == len(files)-1 { //if it's the last element, clean things up
					if len(files) < numFiles {
						files = append(files, newFile)
					} else {
						files = files[:numFiles]
					}

					*smallest = files[len(files)-1].size

					break
				}
				continue
			} else {
				//insert newFile in the correct position
				files = append(files[:i], append([]file{newFile}, files[i:]...)...)

				if len(files) > numFiles {
					files = files[:numFiles]
				}

				*smallest = files[len(files)-1].size

				break
			}
		}
	}

	return files
}

func min(ints []int) int {
	min := ints[0]
	for _, v := range ints {
		if v < min {
			min = v
		}
	}
	return min
}
