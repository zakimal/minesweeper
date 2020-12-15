package storage

type Driver interface {
	New() error
}