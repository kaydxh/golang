package webserver

import "fmt"

func (w *Web) GetBindHostPort() string {
	addr := w.GetBindAddress()

	return fmt.Sprintf("%s:%d:", addr.GetHost(), addr.GetPort())
}
