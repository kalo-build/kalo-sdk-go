package kalo

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// FileStore provides filesystem operations for a store.
type FileStore interface {
	// ReadFile reads the contents of a file.
	ReadFile(path string) ([]byte, error)

	// WriteFile writes data to a file, creating it if necessary.
	WriteFile(path string, data []byte) error

	// ListDir returns the entries in a directory.
	ListDir(path string) ([]DirEntry, error)

	// Stat returns information about a file or directory.
	Stat(path string) (FileInfo, error)

	// Exists returns true if the path exists.
	Exists(path string) bool
}

// DirEntry represents an entry in a directory listing.
type DirEntry struct {
	Name  string
	IsDir bool
}

// FileInfo represents information about a file.
type FileInfo struct {
	Name    string
	Size    int64
	Mode    fs.FileMode
	ModTime time.Time
	IsDir   bool
}

// wasiFileStore implements FileStore using WASI filesystem operations.
type wasiFileStore struct {
	storeName string
	basePath  string
	err       error
}

func (s *wasiFileStore) ReadFile(path string) ([]byte, error) {
	if s.err != nil {
		return nil, s.err
	}
	fullPath := filepath.Join(s.basePath, path)
	return os.ReadFile(fullPath)
}

func (s *wasiFileStore) WriteFile(path string, data []byte) error {
	if s.err != nil {
		return s.err
	}
	fullPath := filepath.Join(s.basePath, path)

	// Ensure parent directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, data, 0644)
}

func (s *wasiFileStore) ListDir(path string) ([]DirEntry, error) {
	if s.err != nil {
		return nil, s.err
	}
	fullPath := filepath.Join(s.basePath, path)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	result := make([]DirEntry, len(entries))
	for i, e := range entries {
		result[i] = DirEntry{
			Name:  e.Name(),
			IsDir: e.IsDir(),
		}
	}
	return result, nil
}

func (s *wasiFileStore) Stat(path string) (FileInfo, error) {
	if s.err != nil {
		return FileInfo{}, s.err
	}
	fullPath := filepath.Join(s.basePath, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
	}, nil
}

func (s *wasiFileStore) Exists(path string) bool {
	if s.err != nil {
		return false
	}
	fullPath := filepath.Join(s.basePath, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

