package option

import (
	"fmt"
	"os"
)

type FileInfo struct {
	Path     string
	FileInfo os.FileInfo
}

func (e *FileInfo) String() string {
	return fmt.Sprintf(
		"{path: [%v], size: [%v], modeTime: [%v]}",
		e.Path,
		e.FileInfo.Size(),
		e.FileInfo.ModTime(),
	)
}

type ExtractMsg struct {
	FileInfo *FileInfo
	Error    error
}

func (e *ExtractMsg) String() string {
	return fmt.Sprintf(
		"{path: [%v], size: [%v], modeTime: [%v]}",
		e.FileInfo.Path,
		e.FileInfo.FileInfo.Size(),
		e.FileInfo.FileInfo.ModTime(),
	)

}
