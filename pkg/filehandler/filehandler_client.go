package filehandler

type Path string

type FileHandler struct {
}

type File struct {
	Path    Path
	Content []byte
}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}
