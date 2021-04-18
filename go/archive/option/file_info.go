package option

import "os"

type FileInfo struct {
	Path     string
	FileInfo os.FileInfo
}

type ExtractMsg struct {
	FileInfo *FileInfo
	Error    error
}
