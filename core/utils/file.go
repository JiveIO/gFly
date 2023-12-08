package utils

import "strings"

// FileExt Extract extension of file.
func FileExt(fileName string) string {
	filePart := strings.Split(fileName, ".")

	return filePart[len(filePart)-1]
}

// FileName Extract file name of file path.
func FileName(filePath string) string {
	filePart := strings.Split(filePath, "/")

	return filePart[len(filePart)-1]
}
