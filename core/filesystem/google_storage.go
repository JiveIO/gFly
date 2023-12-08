package filesystem

import (
	"app/core/errors"
	"app/core/log"
	"os"
	"time"
)

//---------------------------------------
//				Structure				|
//---------------------------------------

// NewGoogleStorage Create Google storage.
func NewGoogleStorage() *GoogleStorage {
	return &GoogleStorage{}
}

type GoogleStorage struct {
}

// ---------------------------------------
//			Implement IStorage			|
// ---------------------------------------

func (s *GoogleStorage) Put(path, contents string, options ...interface{}) bool {
	return true
}

func (s *GoogleStorage) PutFile(path string, fileSource *os.File, options ...interface{}) bool {
	return true
}

func (s *GoogleStorage) Delete(path string) bool {
	return true
}

func (s *GoogleStorage) Copy(from, to string) bool {
	return true
}

func (s *GoogleStorage) Move(from, to string) bool {
	return true
}

func (s *GoogleStorage) Exists(path string) bool {
	return true
}
func (s *GoogleStorage) Get(path string) ([]byte, error) {
	return nil, nil
}

func (s *GoogleStorage) Size(path string) int64 {
	return 1
}

func (s *GoogleStorage) LastModified(path string) time.Time {
	return time.Now()
}

func (s *GoogleStorage) MakeDir(dir string) bool {
	return true
}

func (s *GoogleStorage) DeleteDir(dir string) bool {
	return true
}

func (s *GoogleStorage) Append(path, data string) bool {
	log.Errorf("Unable to append data %s into %s. Here's why: %v\n", path, data, errors.NotYetImplemented.Error())

	return false
}
