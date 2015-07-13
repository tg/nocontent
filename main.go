package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// pathError return error string prefixed with path if it's not already there
func pathError(err error, path string) string {
	s := err.Error()
	if !strings.Contains(s, path) {
		s = fmt.Sprintf("%s: %s", path, s)
	}
	return s
}

func main() {
	var noCounters bool
	flag.BoolVar(&noCounters, "C", false, "Don't print counters, only filenames")
	flag.Parse()

	// Get directories to walk through; if none, use current
	dirs := flag.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	// Use two threads – for walking and reading
	runtime.GOMAXPROCS(2)
	log.SetFlags(0)

	// Channel storing file paths to check
	var files = make(chan string, 100)

	// Spawn file walker. Will close the files channel when done.
	go func() {
		// Ignore walk errors as we never propagate them
		for _, root := range dirs {
			root = filepath.Clean(root) + "/"
			filepath.Walk(root, func(path string, info os.FileInfo, err error) (_ error) {
				if err != nil {
					log.Println(pathError(err, path))
					return
				}
				if !info.Mode().IsRegular() {
					return
				}
				files <- path
				return
			})
		}
		close(files)
	}()

	// Check files
	for path := range files {
		f, err := os.Open(path)
		if err == nil {
			var n int
			n, err = ReadZeros(f)
			if err == nil && n >= 0 {
				if noCounters {
					fmt.Println(path)
				} else {
					fmt.Printf("%10d\t%s\n", n, path)
				}
			}
			f.Close()
		}
		if err != nil && err != errNonZeroByte {
			log.Println(pathError(err, path))
		}
	}
}
