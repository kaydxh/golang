package webserver_test

import (
	"testing"

	"context"

	webserver_ "github.com/kaydxh/golang/pkg/webserver"
)

func TestNew(t *testing.T) {
	config := new(webserver_.Config)
	s, err := config.Complete().New()
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}
	s.InstallWebHandlers()
	prepared, err := s.PrepareRun()
	if err != nil {
		t.Errorf("failed to PrepareRun err: %v", err)
	}

	prepared.Run(context.Background())
}
