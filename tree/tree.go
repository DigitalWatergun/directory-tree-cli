package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/digitalwatergun/directory-tree-cli/config"
)

type entry struct {
	path  string
	depth int
	isDir bool
}

func WalkTree(root string, showFiles bool) ([]string, error) {
	var entries []entry

	isIgnored := func(name string) bool {
		return slices.Contains(config.IgnoreList, name)
	}

	if err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		base := d.Name()
		if d.IsDir() && isIgnored(base) {
			rel, _ := filepath.Rel(root, path)
			depth := 0
			if path != root {
				depth = len(strings.Split(rel, string(os.PathSeparator)))
			}
			entries = append(entries, entry{path, depth, true})
			return filepath.SkipDir
		}

		if path == root {
			entries = append(entries, entry{path, 0, true})
			return nil
		}
		if !showFiles && !d.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		depth := len(filepath.SplitList(rel))
		entries = append(entries, entry{path, depth, d.IsDir()})
		return nil
	}); err != nil {
		return nil, fmt.Errorf("walking %q: %w", root, err)
	}

	var lines []string
	for _, entry := range entries {
		prefix := strings.Repeat("│   ", entry.depth)

		name := filepath.Base(entry.path)
		if entry.isDir {
			name += "/"
		}
		lines = append(lines, prefix+name)
	}
	return lines, nil
}
