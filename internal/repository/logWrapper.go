package repository

import "log"

type LogWrapper struct {
	repo Repository
}

func NewLogWrapper(repo Repository) *LogWrapper {
	return &LogWrapper{
		repo: repo,
	}
}

func (l *LogWrapper) Get(name string) (*Cat, error) {
	log.Printf("Repo: get by name %s\n", name)
	return l.repo.Get(name)
}
func (l *LogWrapper) Delete(name string) error {
	log.Printf("Repo: delete by name %s\n", name)
	return l.repo.Delete(name)
}
func (l *LogWrapper) Insert(cat *Cat) (err error) {
	log.Printf("Repo: insert cat %#v\n", cat)
	return l.repo.Insert(cat)
}
func (l *LogWrapper) Update(catUpdate *CatUpdate) error {
	log.Printf("Repo: update cat %#v\n", catUpdate)
	return l.repo.Update(catUpdate)
}
func (l *LogWrapper) Find(query *Query) ([]*Cat, error) {
	log.Printf("Repo: find by query %#v\n", query)
	return l.repo.Find(query)
}
func (l *LogWrapper) Destroy() error {
	return l.repo.Destroy()
}
