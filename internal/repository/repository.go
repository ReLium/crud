package repository

type Cat struct {
	Name       string
	Gender     string
	Color      string
	Vaccinated bool
}

type CatUpdate struct {
	Name       string
	Gender     string
	Color      string
	Vaccinated *bool
}

type Query struct {
	Gender     string
	Color      string
	Vaccinated *bool
}

//go:generate mockery --dir . --inpackage --name Repository

type Repository interface {
	Get(name string) (*Cat, error)
	Delete(name string) error
	Insert(cat *Cat) (err error)
	Update(catUpdate *CatUpdate) error
	Find(query *Query) ([]*Cat, error)
	Destroy() error
}
