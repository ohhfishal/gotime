package entry

type FileHandler struct {
	filename string
}

func NewFileHandler(filename string) (*FileHandler, error) {
	return &FileHandler{
		filename: filename,
	}, nil
}

func (handler *FileHandler) CreateEntry(newEntry Entry) error {
	return AppendFile(handler.filename, newEntry)
}

func (handler *FileHandler) GetAllEntries() ([]Entry, error) {
	return ReadAllFromFile(handler.filename)
}
