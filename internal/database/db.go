package database

import "context"

//go:generate mockery --inpackage --case snake --disable-version-string --name "EntityStorage"
type EntityStorage interface {
	GetEntities(context.Context) ([]Entity, error)
}

type Entity struct {
	ID          string
	TwitterId   string
	DisplayName string
}
