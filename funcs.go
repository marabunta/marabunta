package marabunta

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
)

// md5sum return md5 checksum of given file
func md5sum(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// inSlice find if item is in slice
func inSlice(s []string, item string) bool {
	for _, i := range s {
		if i == item {
			return true
		}
	}
	return false
}

// isDir return true if path is a dir
func isDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if m := f.Mode(); m.IsDir() && m&400 != 0 {
		return true
	}
	return false
}

// isFile return true if path is a regular file
func isFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if m := f.Mode(); !m.IsDir() && m.IsRegular() && m&400 != 0 {
		return true
	}
	return false
}

// GetHome returns the $HOME/.marabunta
func GetHome() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("error getting user home: %s", err)
		}
		home = usr.HomeDir
	}
	home = filepath.Join(home, ".marabunta")
	if err := os.MkdirAll(home, os.ModePerm); err != nil {
		return "", err
	}
	return home, nil
}
