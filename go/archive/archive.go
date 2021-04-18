package archive

import (
	"github.com/kaydxh/golang/go/archive/option"
)

type Archvier interface {
	Extract(srcFile, destDir string) ([]string, error)
	ExtractStream(srcFile, destDir string) (<-chan option.ExtractMsg, error)
}
