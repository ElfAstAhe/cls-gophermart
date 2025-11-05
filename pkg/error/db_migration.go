package error

import "fmt"

type DBMigrationError struct {
	Migration string
	Err       error
}

var DBMigrationErr *DBMigrationError

func NewDBMigrationError(migration string, err error) *DBMigrationError {
	return &DBMigrationError{
		Migration: migration,
		Err:       err,
	}
}

func (dm *DBMigrationError) Error() string {
	return fmt.Sprintf("migration error with migration [%s]: [%v]", dm.Migration, dm.Err)
}

func (dm *DBMigrationError) Unwrap() error {
	return dm.Err
}
