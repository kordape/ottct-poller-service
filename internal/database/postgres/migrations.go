package postgres

import (
	_ "embed"
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	model "github.com/kordape/ottct-main-service/pkg/db"
	"gorm.io/gorm"
)

var (
	//go:embed seed/202303301900.sql
	seed202303301900 string
)

func (db *DB) Migrate() error {
	m := gormigrate.New(db.db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "entity-schema-202303301900",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&model.Entity{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("entities")
			},
		},
		{
			ID: "insert-entities-202303301900",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec(seed202303301900).Error
			},
		},
	})

	if err := m.Migrate(); err != nil {
		db.log.Error(fmt.Errorf("Could not migrate: %v", err))
		return fmt.Errorf("Migration failed: %v", err)
	}

	db.log.Info("Migration run successfully")

	return nil
}
