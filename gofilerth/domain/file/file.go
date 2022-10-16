package file

type File struct {
	path        string
	isDirectory bool
}

type Dto struct {
	Path        string
	IsDirectory bool
}

func New(path string) *File {
	return &File{path: path}
}

func FromDto(dto Dto) *File {
	return &File{
		path:        dto.Path,
		isDirectory: dto.IsDirectory,
	}
}

func (file *File) IsDirectory() bool {
	return file.isDirectory
}

func (file *File) Path() string {
	return file.path
}
