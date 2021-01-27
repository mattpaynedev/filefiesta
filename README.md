# File Fiesta

---

File Fiesta is a small command line application that outputs the largest files within a given directory.

---

### Background:
File Fiesta is useful for managing the limited capacity of online file storage/sync service (Dropbox/Google Drive/Box/etc.). The program outputs the largest files in a given directory, along with their locations and size. When you are near capacity for your storage service, you can use File Fiesta to locate the largest files in your drive (videos, raw photos, etc.) and decided whether you can delete them to free up space. Note that the program will only search through files that are synced with your machine, so it will not capture files that are only stored in the cloud.


### Installation:


### How to:


### Under the hood:
The File Fiesta app takes user defined values via command line flags and outputs the largest files in the chosen directory, sorted from largest to smallest. By default the program searches the current directory and returns up to 20 of the largest files and ignores hidden folders (whose directory names begin with '.', for example `.git`)

Our `fileSearch` function utilizes the `filepath.Walk` function (from the standard lirbary) to do a recursive, depth-first search through the subject directory. `filepath.Walk` takes a directory address string and `WalkFunc` type function. The `WalkFunc` is an anonymous function that allows for logic to be applied to each node in the depth-first search. 

In this case, we test each node to see whether it should be included in the final output. Additionally, we filter out files/folders we don't want to search and keep track of the number of files seen and the total size of the directory we're searching.

Following the `fileSearch` function, we take the results and print them to the terminal utilizing the `golang.org/x/text/message` package. `message` allows us to format numbers using 000s comma separators, which is not readily doable using the `fmt` package. We then iterate over our resulting array, printing the name, location, and size (in MB) of the largest files.
