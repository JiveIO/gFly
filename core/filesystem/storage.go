package filesystem

import (
	"app/core/utils"
	"os"
	"time"
)

//---------------------------------------
//				Structure				|
//---------------------------------------

type StorageType string

const (
	StorageS3     = StorageType("s3")
	StorageGoogle = StorageType("google")
	StorageLocal  = StorageType("local")
)

// New Create storage instance.
func New(storageType ...StorageType) IStorage {
	var storage StorageType

	storage = StorageType(utils.Getenv("FILESYSTEM_TYPE", "local"))

	if len(storageType) > 0 {
		storage = storageType[0]
	}

	switch storage {
	case StorageS3:
		return NewS3Storage()
	case StorageGoogle:
		return NewGoogleStorage()
	default:
		return NewLocalStorage()
	}
}

//---------------------------------------
//				Interfaces				|
//---------------------------------------

// IStorage Storage interface
type IStorage interface {
	// -- Main actions

	// Put Create a file with content string
	Put(path, contents string, options ...interface{}) bool
	// PutFile Create a file from another file source
	PutFile(path string, fileSource *os.File, options ...interface{}) bool
	// Delete Remove a file
	Delete(path string) bool
	// Copy Clone file to another location
	Copy(from, to string) bool
	// Move Switch file location to new place
	Move(from, to string) bool

	//-- File manipulation

	// Exists Check existed file
	Exists(path string) bool
	// Get Receive a file content
	Get(path string) ([]byte, error)
	// Size Get file size
	Size(path string) int64
	// LastModified Obtains last modified of file
	LastModified(path string) time.Time

	//-- Directory

	// MakeDir Create new directory
	MakeDir(dir string) bool
	// DeleteDir Remove empty directory
	DeleteDir(dir string) bool

	//-- Data manipulation

	// Append Add string content to bottom file
	Append(path, data string) bool
}
