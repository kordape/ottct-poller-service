package database

//go:generate mockery --inpackage --case snake --disable-version-string --name "EntityStorage"
type EntityStorage interface {
	GetEntities() ([]Entity, error)
}

type Entity struct {
	ID          string
	TwitterId   string
	DisplayName string
}
