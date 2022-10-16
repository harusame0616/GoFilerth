package file

type Repository interface {
	GetOneByPath(destinationPath string) (*File, error)
}
