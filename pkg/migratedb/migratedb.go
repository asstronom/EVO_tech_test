package migratedb

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateUp(dburl string) error {
	m, err := migrate.New("file://pkg/migratedb/sql", dburl)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err.Error() != "no change" {
		return err
	}
	return nil
}
