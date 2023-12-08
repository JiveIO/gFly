package filesystem

import (
	"app/core/log"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalStoragePut(t *testing.T) {
	fsDefault := NewLocalStorage()

	fsDefault.MakeDir("logs")
	fsDefault.MakeDir("tmp")

	filePath := "logs/logs.log"

	// Put content
	fsDefault.Put(filePath, "Hello\n")
}

func TestLocalStoragePutFile(t *testing.T) {
	fsDefault := NewLocalStorage()

	path := "storage/logs/logs.log"
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		log.Infof("Couldn't open file %v. Here's why: %v\n", path, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Unable to close file %q, %v", path, err)
		}
	}(file)

	// Put file
	fsDefault.PutFile("logs/logs-file.log", file)
}

func TestLocalStorageCopy(t *testing.T) {
	fsDefault := NewLocalStorage()

	from := "logs/logs.log"
	to := "tmp/logs-file.log"

	// Copy file
	fsDefault.Copy(from, to)
}

func TestLocalStorageMove(t *testing.T) {
	fsDefault := NewLocalStorage()

	from := "logs/logs.log"
	to := "tmp/logs-file1.log"
	to2 := "tmp/logs-file2.log"
	to3 := "logs/logs-file3.log"

	// Copy file
	fsDefault.Move(from, to)

	// Move file
	fsDefault.Move(to, to2)

	// Move file
	fsDefault.Move(to2, to3)
}

func TestLocalStorageExists(t *testing.T) {
	fsDefault := NewLocalStorage()

	path := "logs/logs-file3.log"

	fsDefault.Exists(path)
}

func TestLocalStorageGet(t *testing.T) {
	fsDefault := NewLocalStorage()
	path := "logs/logs-file3.log"

	_, err := fsDefault.Get(path)
	if err != nil {
		fmt.Printf("Error %v", err)
	}
}

func TestLocalStorageInfo(t *testing.T) {
	fsDefault := NewLocalStorage()
	path := "logs/logs-file3.log"

	_ = fsDefault.Size(path)
	_ = fsDefault.LastModified(path)
}

func TestLocalStorageDir(t *testing.T) {
	fsDefault := NewLocalStorage()

	fsDefault.MakeDir("tmp_one/two")
	fsDefault.DeleteDir("tmp_one/two")
	fsDefault.DeleteDir("tmp_one")
	fsDefault.Append("logs/logs-file3.log", "Hello world new line\n")
}

func TestLocalStorageDelete(t *testing.T) {
	fsDefault := NewLocalStorage()
	files := []string{
		"logs/logs-file.log",
		"logs/logs-file3.log",
		"tmp/logs-file.log",
	}

	for _, filePath := range files {
		fsDefault.Delete(filePath)
	}
	fsDefault.DeleteDir("logs")
	fsDefault.DeleteDir("tmp")
}
