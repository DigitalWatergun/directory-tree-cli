package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/digitalwatergun/directory-tree-cli/tree"
)

func main() {
	showFiles := flag.Bool("files", false, "Show files in addition to directories")
	out := flag.String("out", "", "Path to output .txt file. (Default: stdout)")
	flag.Parse()

	root := "."

	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	var outPath string
	if *out != "" {
		outPath = *out
	} else {
		var base string
		if root == "." {
			wd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				os.Exit(1)
			}
			base = filepath.Base(wd)
		} else {
			base = filepath.Base(root)
		}
		outPath = base + ".txt"
	}

	lines, err := tree.WalkTree(root, *showFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var file *os.File
	if outPath != "" {
		file, err = os.Create(outPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output %q: %v\n", outPath, err)
			os.Exit(1)
		}
		defer func() {
			if cerr := file.Close(); cerr != nil {
				fmt.Fprintf(os.Stderr, "Error closing file %q: %v\n", outPath, cerr)
			}
		}()
	} else {
		file = os.Stdout
	}

	for _, line := range lines {
		if _, err := fmt.Fprintln(file, line); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing line %q: %v\n", line, err)
			os.Exit(1)
		}
	}

	fmt.Fprintf(os.Stderr, "Wrote Tree to %s\n", outPath)
}
