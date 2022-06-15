# golang
Golang lib is a programming toolkit for building microservices in Go. It has very useful interface or service to develop application.
# To start using golang
To make use web server of golang, also golang has very useful libs, It speeds up your development.
```
import (
	"testing"

	"context"

	viper_ "github.com/kaydxh/golang/pkg/viper"
	webserver_ "github.com/kaydxh/golang/pkg/webserver"
)

func main() {
	cfgFile := "./webserver.yaml"
	config := webserver_.NewConfig(webserver_.WithViper(viper_.GetViper(cfgFile, "web")))

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
```

# Evolution
Golang started in Oct 8, 2020. 

# Contributing
If you need support, start with your branch, and create a pull request for us. We appreciate your help!
