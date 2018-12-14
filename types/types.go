package types

import "io"

// Manager - interface for Managers to follow
type Manager interface {
	SetData([]Data)
	LoadDataFromReader(io.Reader) ([]Data, error)
	Data() []Data
	ManagerValidator
}

// Data - is the underlying data struct
type Data interface {
	Valid() []error
}

type ErroredRecord struct {
	Err  []error
	Data Data
}

// ManagerValidator - This is responsible for the thing that validates the data in the manager
type ManagerValidator interface {
	ValidateCollection() []ErroredRecord
}
