package filehasher

import (
	"context"
	"go-hasher/internal/controllers/hasher"
	"go-hasher/pkg/appcontext"
	"go-hasher/pkg/filehandler"
)

type FileHasherController struct {
	fileHandler *filehandler.FileHandler
	hasher      *hasher.HasherController
}

func NewFileHasherController() *FileHasherController {
	return &FileHasherController{
		fileHandler: filehandler.NewFileHandler(),
		hasher:      hasher.NewHasherController(),
	}
}

func (fhc *FileHasherController) HashFile(ctx context.Context, path filehandler.Path) (string, error) {
	valid, _ := fhc.fileHandler.ValidateFile(path)
	if !valid {
		return "", nil
	}
	memoryCache := appcontext.GetMemoryCache(ctx)
	if memoryCache == nil {
		return "", nil
	}
	if hash, exists := memoryCache.Get(string(path)); exists {
		return hash, nil
	}

	file, err := fhc.fileHandler.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := fhc.hasher.HashFile(file.Content)
	memoryCache.Set(string(path), hash)
	return hash, nil
}
