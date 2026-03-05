package util

import (
	"os"
	"path/filepath"
	"strings"
)

func NormalizePathArgs(raw []string) []string {
	out := make([]string, 0, len(raw))
	for _, entry := range raw {
		trimmed := strings.TrimSpace(entry)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}

func ExpandUploadPaths(paths []string) ([]string, error) {
	result := make([]string, 0)
	seen := map[string]struct{}{}

	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			abs, err := filepath.Abs(p)
			if err != nil {
				return nil, err
			}
			if _, ok := seen[abs]; !ok {
				seen[abs] = struct{}{}
				result = append(result, abs)
			}
			continue
		}

		err = filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			abs, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			if _, ok := seen[abs]; ok {
				return nil
			}
			seen[abs] = struct{}{}
			result = append(result, abs)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
