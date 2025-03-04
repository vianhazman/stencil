package models

import "errors"

var (
	//ErrSnapshotNotFound used when snapshot is not found
	ErrSnapshotNotFound = errors.New("not found")
)

// Snapshot represents specific version of protodescriptorset
type Snapshot struct {
	ID        int64  `binding:"required"`
	Namespace string `binding:"required"`
	Name      string `binding:"required"`
	Version   string `binding:"required,version"`
	Latest    bool   `binding:"required"`
}
