package filer

type Repository interface {
	GetOneById(id string) (*Filer, error)
	Save(filer *Filer) error
}
