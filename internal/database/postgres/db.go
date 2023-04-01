package postgres

import (
	"errors"
	"fmt"

	"github.com/kordape/ottct-poller-service/internal/database"
	"github.com/kordape/ottct-poller-service/pkg/logger"
	"gorm.io/gorm"
)

var _ database.EntityStorage = &DB{}

type DB struct {
	db  *gorm.DB
	log *logger.Logger
}

func New(db *gorm.DB, log *logger.Logger) (*DB, error) {
	d := &DB{
		db:  db,
		log: log,
	}

	err := d.validate()

	if err != nil {
		return nil, fmt.Errorf("Failed to initialize DB: %v", err)
	}

	return d, nil
}

func (db *DB) validate() error {
	if db.db == nil {
		return errors.New("validation error: db is nil")
	}

	if db.log == nil {
		return errors.New("validation error: log is nil")
	}

	return nil
}
