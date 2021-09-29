package mount_test

import (
	"testing"

	mount_ "github.com/kaydxh/golang/pkg/storage/mount"
)

func TestMountCeph(t *testing.T) {
	mountPoint, err := mount_.MountCeph("/mnt", "admin", "127.0.0.1:9090", "password", "/data", false, 30)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	t.Logf("mountPoint: %v", mountPoint)
}
