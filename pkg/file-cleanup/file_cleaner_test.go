package filecleanup_test

import (
	"testing"
	"time"

	filecleanup_ "github.com/kaydxh/golang/pkg/file-cleanup"
)

func TestFileCleanupWithMaxAge(t *testing.T) {

	err := filecleanup_.FileCleanup(
		"/Users/kayxhding/workspace/studyspace/git-kayxhding/github.com/kaydxh/golang/pkg/file-cleanup/*",
		filecleanup_.WithMaxAge(time.Second),
	)
	if err != nil {
		t.Fatalf("faild to FileCleanup, err: %v", err)
	}
}
