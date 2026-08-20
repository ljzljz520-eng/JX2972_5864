package store

import (
	"os"
	"path/filepath"
)

func EnsureParent(path string) error {
	parent := filepath.Dir(path)
	if parent == "." || parent == "" {
		return nil
	}
	return os.MkdirAll(parent, 0755)
}

func Exists(path string) bool { _, err := os.Stat(path); return err == nil }

func Remove(path string) error {
	if !Exists(path) {
		return nil
	}
	return os.Remove(path)
}

func CopyFile(source, target string) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	if err := EnsureParent(target); err != nil {
		return err
	}
	return os.WriteFile(target, data, 0600)
}
