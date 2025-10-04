package filehandler

import (
	"os"
)

func (f *FileHandler) ReadFile(path Path) (File, error) {
	content, err := os.ReadFile(string(path))
	if err != nil {
		return File{}, err
	}
	return File{Path: path, Content: content}, nil
}

func (f *FileHandler) WriteFile(file File) error {
	return os.WriteFile(string(file.Path), file.Content, 0644)
}

func (f *FileHandler) ValidateFile(path Path) (bool, error) {
	info, err := os.Stat(string(path))
	if err != nil {
		return false, err
	}
	if info.IsDir() {
		return false, nil
	}
	return true, nil
}
