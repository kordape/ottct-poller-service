package postgres

import (
	"fmt"

	model "github.com/kordape/ottct-main-service/pkg/db"
	"github.com/kordape/ottct-poller-service/internal/database"
)

func (db *DB) GetEntities() ([]database.Entity, error) {
	var persistentEntities []model.Entity
	err := db.db.Find(&persistentEntities).Error
	if err != nil {
		return nil, fmt.Errorf("Error getting entities from db: %w", err)
	}

	entities := make([]database.Entity, len(persistentEntities))
	for i, e := range persistentEntities {
		entities[i] = database.Entity{
			ID:          e.ID,
			TwitterId:   e.TwitterId,
			DisplayName: e.DisplayName,
		}
	}

	return entities, nil
}
